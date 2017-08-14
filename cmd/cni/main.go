package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ns"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-cni/pkg/utils"

	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/types"
	"net"
	"os"
	"runtime"
	"time"
	"github.com/yunify/qingcloud-cni/pkg/nicmanagr"

	"github.com/containernetworking/plugins/pkg/ipam"
)

// PluginConf is whatever you expect your configuration json to be. This is whatever
// is passed in on stdin. Your plugin may wish to expose its functionality via
// runtime args, see CONVENTIONS.md in the CNI spec.
type PluginConf struct {
	types.NetConf // You may wish to not nest this type
	Args          struct {
		BindAddr string `json:"bindaddr,omitempty"`
	} `json:"args,omitempty"`
}

const (
	DefaultDeleteTimeout = 30 * time.Second
)

// parseConfig parses the supplied configuration (and prevResult) from stdin.
func parseConfig(stdin []byte) (*PluginConf, error) {
	conf := PluginConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}
	return &conf, nil
}

func cmdAdd(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:0", remote.WithEndpointWriterBatchSize(10000))
	manager := actor.NewPID(conf.Args.BindAddr, nicmanagr.NicManagerActorName)
	result, err := manager.RequestFuture(&messages.AllocateNicMessage{
		Name: args.ContainerID,
	}, nicmanagr.DefaultCreateNicTimeout).Result()
	if err != nil {
		log.Errorf("Failed to get nic: %v", err)
		return err
	}
	response := result.(*messages.AllocateNicReplyMessage)
	if response.EndpointID == "" {
		log.Errorf("Failed to get nic: %v")
		return fmt.Errorf("Failed to get nic: %v")
	}

	gatewaymgr := actor.NewPID(conf.Args.BindAddr,nicmanagr.GatewayManagerActorName)
	result,err = gatewaymgr.RequestFuture(&messages.GetGatewayMessage{
		Vxnet: response.NetworkID,
	},nicmanagr.DefaultGetGatewayTimeout).Result()
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		log.Errorf("Failed to get default gateway.%v",err)
		return err
	}
	gateway := result.(*messages.GetGatewayReplyMessage)
	if gateway.NetworkCIDR == ""{
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		log.Errorf("Failed to get default gateway. response is empty")
		return err
	}

	//bind iface to ns
	iface, err := utils.LinkByMacAddr(response.EndpointID)
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		return fmt.Errorf("LinkByMacAddr err %s, delete Nic %s", err.Error(), response.EndpointID)
	}

	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		return fmt.Errorf("Failed to get network namespace")
	}
	defer netns.Close()

	if err = netlink.LinkSetNsFd(iface, int(netns.Fd())); err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		return fmt.Errorf("LinkSetNsFd err %s, delete Nic %s", err.Error(), response.EndpointID)
	}

	ifacename := iface.Attrs().Name

	//generate result
	netiface := &current.Interface{Name: args.IfName, Mac: response.EndpointID, Sandbox: args.ContainerID}
	nicIP:=net.ParseIP(response.Address)
	_,ipNet,_:=net.ParseCIDR(gateway.NetworkCIDR)
	ipConfig := &current.IPConfig{
		Address: net.IPNet{
			IP: nicIP,
			Mask: ipNet.Mask,
		},
		Version: "4",
		Gateway: net.ParseIP(gateway.Gateway),
	}
	res := &current.Result{
		Interfaces: []*current.Interface{netiface},
		IPs: []*current.IPConfig{ipConfig},
		Routes:[]*types.Route{
			&types.Route{Dst: net.IPNet{IP: net.IPv4zero, Mask: net.IPv4Mask(0,0,0,0)}, GW: net.ParseIP(gateway.Gateway)},
		},
	}

	//configure nic
	err = netns.Do(func(_ ns.NetNS) error {
		nsiface, err := netlink.LinkByName(ifacename)
		if err != nil {
			return fmt.Errorf("failed to get link by name %q: %v", ifacename, err)
		}

		if err := netlink.LinkSetDown(nsiface); err != nil {
			return fmt.Errorf("failed to set %q Down: %v", args.IfName, err)
		}

		if err := netlink.LinkSetName(nsiface, args.IfName); err != nil {
			return fmt.Errorf("failed to setname %q: %v", args.IfName, err)
		}

		if err:= ipam.ConfigureIface(args.IfName,res); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		return fmt.Errorf("Failed to configure interface %s, delete Nic %s", err.Error(), response.EndpointID)
	}



	return types.PrintResult(res, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}
	remote.Start("127.0.0.1:0", remote.WithEndpointWriterBatchSize(10000))
	manager := actor.NewPID(conf.Args.BindAddr, nicmanagr.NicManagerActorName)

	err = ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		ifName := args.IfName
		iface, err := netlink.LinkByName(ifName)
		if err != nil {
			return fmt.Errorf("failed to lookup %q: %v", ifName, err)
		}
		manager.RequestFuture(&messages.DeleteNicMessage{Nicname: args.ContainerID,Nicid:iface.Attrs().HardwareAddr.String()}, DefaultDeleteTimeout).Wait()

		if err := netlink.LinkSetDown(iface); err != nil {
			return fmt.Errorf("Failed to set link down:%v", err)
		}
		if err = netlink.LinkSetNsPid(iface, os.Getpid()); err != nil {
			return fmt.Errorf("Failed to set namespace to default ns, %v", err)
		}
		return nil
	})

	return err
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

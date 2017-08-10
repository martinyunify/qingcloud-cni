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

	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/types"
	"net"
	"os"
	"runtime"
	"time"
	"syscall"
)

// PluginConf is whatever you expect your configuration json to be. This is whatever
// is passed in on stdin. Your plugin may wish to expose its functionality via
// runtime args, see CONVENTIONS.md in the CNI spec.
type PluginConf struct {
	types.NetConf // You may wish to not nest this type
	Args struct {
		BindAddr string `json:"bindaddr,omitempty"`
		ActorID  string `json:"actorid,omitempty"`
	} `json:"args,omitempty"`
}

const (
	DefaultDeleteTimeout = 30*time.Second
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
	manager := actor.NewPID(conf.Args.BindAddr, conf.Args.ActorID)
	result, err := manager.RequestFuture(&messages.AllocateNicMessage{
		Name: args.ContainerID,
	}, 60*time.Second).Result()
	if err != nil {
		log.Errorf("Failed to get nic: %v", err)
		return err
	}
	response := result.(*messages.AllocateNicReplyMessage)
	if response.EndpointID == "" {
		log.Errorf("Failed to get nic: %v")
		return fmt.Errorf("Failed to get nic: %v")
	}

	//bind iface to ns
	iface, err := LinkByMacAddr(response.EndpointID)
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		return fmt.Errorf("LinkByMacAddr err %s, delete Nic %s", err.Error(), response.EndpointID)
	}

	netns,err:=ns.GetNS(args.Netns)
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

	//get Default gateway

	//configure nic
	err = netns.Do(func(_ ns.NetNS) error {
		nsiface , err := netlink.LinkByName(ifacename)
		if err != nil {
			return fmt.Errorf("failed to get link by name %q: %v", ifacename, err)
		}

		if err := netlink.LinkSetDown(nsiface); err != nil {
			return fmt.Errorf("failed to set %q Down: %v", args.IfName, err)
		}

		if err := netlink.LinkSetName(nsiface, args.IfName); err != nil {
			return fmt.Errorf("failed to setname %q: %v", args.IfName, err)
		}

		// clean up addr conf
		addrs, err := netlink.AddrList(nsiface, syscall.AF_INET)
		if err == nil && len(addrs) > 0 {
			for _, addr := range addrs {
				err := netlink.AddrDel(iface, &addr)
				if err != nil {
					return fmt.Errorf("AddrDel err %s addr:%+v, Nic %s", err.Error(), addr, iface.Attrs().HardwareAddr)
				}
			}
		}

		// add ip address



		if err := netlink.LinkSetUp(nsiface); err != nil {
			return fmt.Errorf("failed to set %q UP: %v", args.IfName, err)
		}
		return nil
	})
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, DefaultDeleteTimeout).Wait()
		return fmt.Errorf("Failed to configure interface %s, delete Nic %s", err.Error(), response.EndpointID)
	}

	//generate result
	netiface := &current.Interface{Name: args.IfName, Mac: response.EndpointID, Sandbox: args.ContainerID}
	ipConfig := &current.IPConfig{Address: net.IPNet{IP: net.ParseIP(response.Address)}}
	res := &current.Result{Interfaces: []*current.Interface{netiface}, IPs: []*current.IPConfig{ipConfig}}

	return types.PrintResult(res, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}
	remote.Start("127.0.0.1:0", remote.WithEndpointWriterBatchSize(10000))
	manager := actor.NewPID(conf.Args.BindAddr, conf.Args.ActorID)

	err = ns.WithNetNSPath(args.Netns, func(_ ns.NetNS) error {
		ifName := args.IfName
		iface, err := netlink.LinkByName(ifName)
		if err != nil {
			return fmt.Errorf("failed to lookup %q: %v", ifName, err)
		}
		if err := netlink.LinkSetDown(iface); err != nil {
			return fmt.Errorf("Failed to set link down:%v", err)
		}
		if err = netlink.LinkSetNsPid(iface,os.Getpid()); err != nil {
			return fmt.Errorf("Failed to set namespace to default ns, %v",err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = manager.RequestFuture(&messages.DeleteNicMessage{Nicname: args.ContainerID}, DefaultDeleteTimeout).Wait()
	if err != nil {
		log.Errorf("Failed to delete nic from remote :%v", err)
	}
	return err
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

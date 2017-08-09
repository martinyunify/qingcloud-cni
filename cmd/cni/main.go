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
	}, 5*time.Second).Result()
	if err != nil {
		log.Errorf("Failed to get nic: %v", err)
		return err
	}
	response := result.(*messages.AllocateNicReplyMessage)

	//bind iface to ns
	iface, err := LinkByMacAddr(response.EndpointID)
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, 5*time.Second).Wait()
		return fmt.Errorf("LinkSetNsFd err %s, delete Nic %s", err.Error(), response.EndpointID)

	}

	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, 5*time.Second).Wait()
		return fmt.Errorf("failed to open netns %q: %v", args.Netns, err)
	}
	defer netns.Close()

	if err = netlink.LinkSetNsFd(iface, int(netns.Fd())); err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, 5*time.Second).Wait()
		return fmt.Errorf("LinkSetNsFd err %s, delete Nic %s", err.Error(), response.EndpointID)
	}

	if err = netlink.LinkSetUp(iface); err != nil {
		manager.RequestFuture(&messages.DeleteNicMessage{Nicid: response.EndpointID}, 5*time.Second).Wait()
		return fmt.Errorf("failed to set %q UP: %v", iface.Attrs().Name, err)
	}

	//return result
	netiface := &current.Interface{Name: args.IfName, Mac: response.EndpointID, Sandbox: args.ContainerID}
	ipConfig := &current.IPConfig{Address: net.IPNet{IP: net.ParseIP(response.Address)}}
	res := &current.Result{Interfaces: []*current.Interface{netiface}, IPs: []*current.IPConfig{ipConfig}}
	if err != nil {
		log.Errorf("Failed to get previous result,%v", err)
	}
	return types.PrintResult(res, conf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}

	remote.Start("127.0.0.1:0", remote.WithEndpointWriterBatchSize(10000))
	manager := actor.NewPID(conf.Args.BindAddr, conf.Args.ActorID)

	netns, err := ns.GetNS(args.Netns)
	if err != nil {
		return err
	}
	defer netns.Close()
	ifName := args.IfName
	iface, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("failed to lookup %q: %v", ifName, err)
	}
	if err := netlink.LinkDel(iface); err != nil {
		return fmt.Errorf("Failed to del link :%v", err)
	}

	result, err := manager.RequestFuture(&messages.DeleteNicMessage{Nicid: iface.Attrs().HardwareAddr.String()}, 5*time.Second).Result()
	if err != nil {
		log.Errorf("Got error %v", err)
	}
	log.Debug(result)
	return err
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)
}

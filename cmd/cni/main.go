package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"

	"os"
	"time"
)

func cmdDel(args *skel.CmdArgs) error {
	remote.Start("127.0.0.1:0")
	manager := actor.NewPID("127.0.0.1:31080", "nicmanager")
	manager.Tell(messages.DeleteNicMessage{})
	return nil
}

func cmdAdd(args *skel.CmdArgs) error {
	remote.Start("127.0.0.1:0")
	manager := actor.NewPID("127.0.0.1:31080", "nicmanager")
	result, err := manager.RequestFuture(messages.AllocateNicMessage{
		Name: "test",
	}, 30*time.Second).Result()
	if err != nil {
		log.Errorf("Got error %v", err)
	}
	response := result.(*messages.AllocateNicReplyMessage)
	log.Debugf("Got response %s", response)
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}

func init() {
	log.SetOutput(os.Stderr)
}

package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	log "github.com/sirupsen/logrus"

	"time"

	"os"
)

func cmdDel(args *skel.CmdArgs) error {
	manager := actor.NewPID("127.0.0.1:31080", "nicmanager")
	manager.Tell(messages.DeleteNicMessage{})
	return nil
}

func cmdAdd(args *skel.CmdArgs) error {
	manager := actor.NewPID("127.0.0.1:31080", "nicmanager")
	result, err := manager.RequestFuture(&messages.AllocateNicMessage{}, 30*time.Second).Result()
	if err != nil {
		log.Errorf("Got error %v", err)
	}
	response := result.(*messages.AllocateNicReplyMessage)
	log.Debugf("Got response %s", response)
	return nil
}

func main(){
	skel.PluginMain(cmdAdd, cmdDel, version.All)
}

func init() {
	log.SetOutput(os.Stderr)
}
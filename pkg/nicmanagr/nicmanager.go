package nicmanagr

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-cni/pkg/qingcloud"
)

type NicManager struct {
	iface    string
	vxnet    []string
	qingStub *actor.PID
}

func (manager *NicManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case messages.QingcloudInitializeMessage:
		props := actor.FromProducer(qingcloud.NewMockQingCloud)
		var err error
		manager.qingStub, err = actor.SpawnNamed(props, "qingcloud")
		if err != nil {
			log.Errorf("failed to spawn qingcloud actor: %v", err)
		}
		manager.qingStub.Tell(msg)
		ctx.PushBehavior(manager.ProcessMsg)
		log.Debugf("resource stub is activated")
	}
}

func (manager *NicManager) ProcessMsg(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case messages.AllocateNicMessage:
		manager.qingStub.Request(msg, ctx.Sender())
	case *actor.Stopping:
		manager.qingStub.GracefulStop()
		ctx.PopBehavior()
	}
}

func NewNicManager() actor.Actor {
	return &NicManager{}
}

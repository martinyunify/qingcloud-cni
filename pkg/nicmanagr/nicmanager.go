package nicmanagr

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-cni/pkg/qingcloud"
	"time"
)

type NicManager struct {
	qingStub *actor.PID
}

func (manager *NicManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case messages.QingcloudInitializeMessage:
		props := actor.FromProducer(qingcloud.NewMockQingCloud)
		manager.qingStub = actor.Spawn(props)
		manager.qingStub.Tell(msg)
		ctx.PushBehavior(manager.ProcessMsg)
		log.Debugf("resource stub is activated")
	}
}

func (manager *NicManager) ProcessMsg(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case messages.CreateNewNicMessage:
		result, _ := manager.qingStub.RequestFuture(msg, 30*time.Second).Result()
		ctx.Respond(result)
	case *actor.Stopping:
		ctx.PopBehavior()
	}
}

func NewNicManager() actor.Actor {
	return &NicManager{}
}

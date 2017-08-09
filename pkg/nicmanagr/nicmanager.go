package nicmanagr

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-cni/pkg/qingactor"
	"time"
)

type NicManager struct {
	iface    string
	vxnet    []string
	qingStub *actor.PID
}

func (manager *NicManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case messages.QingcloudInitializeMessage:
		props := actor.FromProducer(qingactor.NewQingCloudActor)
		var err error
		manager.qingStub, err = actor.SpawnNamed(props, "qingcloud")
		if err != nil {
			log.Errorf("failed to spawn qingcloud actor: %v", err)
		}
		manager.vxnet = msg.Vxnet
		manager.qingStub.Tell(msg)

		ctx.PushBehavior(manager.ProcessMsg)
		log.Debugf("resource stub is activated")
	}
}

func (manager *NicManager) ProcessMsg(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.AllocateNicMessage:
		result, err := manager.qingStub.RequestFuture(qingactor.CreateNicMessage{
			NetworkID: manager.vxnet[0],
			Nicname:   msg.Name,
		}, 30*time.Second).Result()
		if err != nil {
			log.Errorf("Failed to create new nic: %v", err)
		}
		response := result.(qingactor.CreateNicReplyMessage)
		ctx.Respond(&messages.AllocateNicReplyMessage{
			Name:       msg.Name,
			EndpointID: response.Nic.EndpointID,
			NetworkID:  response.Nic.NetworkID,
			Address:    response.Nic.Address,
		})
		log.Debugf("Allocate new interface.%s", response.Nic.EndpointID)
	case *messages.DeleteNicMessage:
		result, err := manager.qingStub.RequestFuture(qingactor.DeleteNicMessage{
			Nic: msg.Nicid,
		}, 30*time.Second).Result()
		if err != nil || result != nil {
			log.Errorf("Failed to delete nic: %v,%v", err, result)
		}
		ctx.Respond(result)
		log.Debugf("Delete nic.%s", msg.Nicid)
	case *actor.Stopping:
		manager.qingStub.GracefulStop()
		ctx.PopBehavior()
	}
}

func NewNicManager() actor.Actor {
	return &NicManager{}
}

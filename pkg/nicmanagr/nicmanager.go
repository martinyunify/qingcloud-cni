package nicmanagr

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-cni/pkg/qingactor"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/mailbox"
)

const NicManagerActorName = "NicManager"
func init(){
	props := actor.FromProducer(NewNicManager).WithMailbox(mailbox.Bounded(1000))
	_, err := actor.SpawnNamed(props, NicManagerActorName)
	if err != nil {
		log.Error(err)
	} else {
		log.Debugf("Nic Manager is spawned")
	}
}
type NicManager struct {
	iface    string
	vxnet    []string
	qingStub *actor.PID
	policy *CreationPolicy
}

func (manager *NicManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case ResourcePoolInitMessage:
		manager.vxnet = msg.Vxnet
		var err error
		manager.policy,err = NewCreationPolicy(len(manager.vxnet),msg.Policy)
		if err != nil {
			log.Error(fmt.Errorf("Failed to set creation policy %v",err))
			return
		}
		manager.qingStub = actor.NewLocalPID(qingactor.QingCloudActorName)
		ctx.PushBehavior(manager.ProcessMsg)
		log.Debugf("Nicmanager is activated")
	}
}

func (manager *NicManager) ProcessMsg(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.AllocateNicMessage:
		response,err:= manager.createNewNic(msg,1)
		if err != nil || response.Err != nil {
			ctx.Respond(&messages.AllocateNicReplyMessage{})
			return
		}
		ctx.Respond(&messages.AllocateNicReplyMessage{
			Name:       msg.Name,
			EndpointID: response.Nic.EndpointID,
			NetworkID:  response.Nic.NetworkID,
			Address:    response.Nic.Address,
		})
		log.Debugf("Allocate new interface.%s", response.Nic.EndpointID)
	case *messages.DeleteNicMessage:
		result, err := manager.qingStub.RequestFuture(qingactor.DeleteNicMessage{
			Nic:     msg.Nicid,
		}, qingactor.DefaultDeleteNicTimeout).Result()
		reply := result.(qingactor.DeleteNicReplyMessage)
		if err != nil || reply.Err != nil {
			log.Errorf("Failed to delete nic: %v,%v", err, result)
		}
		ctx.Respond(msg)
		log.Debugf("Delete nic.%s", msg.Nicid)
	case *actor.Stopping:
		ctx.PopBehavior()
	}
}

func NewNicManager() actor.Actor {
	return &NicManager{}
}


func (manager *NicManager) createNewNic(msg *messages.AllocateNicMessage,retry int)(*qingactor.CreateNicReplyMessage,error){
	var err error
	for i := 0 ;i < retry; i ++ {
		var result interface{}
		result, err = manager.qingStub.RequestFuture(qingactor.CreateNicMessage{
			NetworkID: manager.vxnet[manager.policy.GetNextItem()],
			Nicname: msg.Name,
		}, qingactor.DefaultCreateNicTimeout).Result()
		response := result.(qingactor.CreateNicReplyMessage)

		if err != nil || response.Err != nil {
			manager.qingStub.RequestFuture(qingactor.DeleteNicMessage{
				Nic: response.Nic.EndpointID,
			}, qingactor.DefaultDeleteNicTimeout).Wait()
			err = fmt.Errorf("Failed to create new nic: %v,%v", err, response.Err)
			log.Error(err)
			manager.policy.UpdateResult(err)
			continue
		}

		manager.policy.UpdateResult(err)
		return &response,err
	}

	return nil,err
}
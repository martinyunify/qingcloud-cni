package nicmanagr

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-cni/pkg/qingactor"
	"github.com/yunify/qingcloud-cni/pkg/utils"
	"time"
)

const NicManagerActorName = "NicManager"

func init() {
	props := actor.FromProducer(NewNicManager).WithMailbox(mailbox.Bounded(1000))
	_, err := actor.SpawnNamed(props, NicManagerActorName)
	if err != nil {
		log.Error(err)
	} else {
		log.Debugf("Nic Manager is spawned")
	}
}

const (
	DefaultCreateNicTimeout = 60 * time.Second
	DefaultDeleteNicTimeout = 60 * time.Second
)

type NicManager struct {
	iface    string
	vxnet    []string
	qingStub *actor.PID
	policy   *CreationPolicy
	cnicache *Nicqueue
	cnicacheSize int
	cachens  string

}

func (manager *NicManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case ResourcePoolInitMessage:
		manager.vxnet = msg.Vxnet
		var err error
		manager.policy, err = NewCreationPolicy(len(manager.vxnet), msg.Policy)
		if err != nil {
			log.Error(fmt.Errorf("Failed to set creation policy %v", err))
			return
		}
		manager.qingStub = actor.NewLocalPID(qingactor.QingCloudActorName)
		manager.cachens = msg.NicNameinCache
		manager.cnicache = &Nicqueue{}
		manager.cnicacheSize = msg.NicCacheSize
		manager.createNewNicInVxnetAsync(manager.cachens, manager.vxnet[manager.policy.GetNextItem()],manager.cnicacheSize)
	case qingactor.CreateNicReplyMessage:

		manager.policy.UpdateResult(msg.Err)
		unusedlist:=manager.cnicache.Enqueue(msg.Nic ...)
		if size:=manager.cnicache.GetPoolShortFall();size >0{
			manager.createNewNicInVxnetAsync(manager.cachens, manager.vxnet[manager.policy.GetNextItem()],size)
			log.Debugf("Nicpool allocate %d more nics", size)
			return
		}
		if len(unusedlist) >0 {
			manager.qingStub.Tell(qingactor.DeleteNicMessage{Nic: unusedlist})
			log.Debugf("Deleted unused %d nic",len(unusedlist))
		}
		ctx.PushBehavior(manager.ProcessMsg)
		log.Debugf("Nicmanager is activated")
	}
}

func (manager *NicManager) ProcessMsg(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.AllocateNicMessage:
		reply := &messages.AllocateNicReplyMessage{}
		if manager.cnicache.IsEmpty() {
			response, err := manager.createOneNewNic(msg.Name, 1)
			if err != nil {
				log.Errorf("Failed to allocate new nic %v",err)
			}else if response.Err != nil {
				log.Errorf("Failed to allocate new nic %v",response.Err)
			} else {
				reply.Name = msg.Name
				reply.EndpointID = response.Nic[0].EndpointID
				reply.NetworkID = response.Nic[0].NetworkID
				reply.Address = response.Nic[0].Address

				//register this nic to pool
				manager.cnicache.AddNewEntries(response.Nic[0])
			}

		} else {
			nic := manager.cnicache.Dequeue()
			reply.Name= msg.Name
			reply.EndpointID = nic.EndpointID
			reply.NetworkID = nic.NetworkID
			reply.Address = nic.Address
			manager.qingStub.Tell(qingactor.ModifyNicNameMessage{Nicid:nic.EndpointID,Nicname:msg.Name})
			manager.createNewNicInVxnetAsync(manager.cachens,manager.vxnet[manager.policy.GetNextItem()],1)
		}
		ctx.Respond(reply)
		log.Debugf("Allocate new interface.%v",reply)
	case *messages.DeleteNicMessage:
		nic := manager.cnicache.GetEntry(&msg.Nicid)
		if manager.cnicache.GetPoolShortFall() > 0 {
			if nic == nil {
				result,err:=manager.qingStub.RequestFuture(qingactor.DescribeNicMessage{Nicid:[]*string{&nic.EndpointID}},qingactor.DefaultQueryNicTimeout).Result()
				if err != nil {
					log.Errorf("Failed to describe recycled nic:%v",err)
				} else if response := result.(qingactor.DescribeNicReplyMessage);response.Err != nil {
					log.Errorf("Failed to describe recycled nic:%v",response.Err)
				} else {
					nic = response.Endpoints[0]
				}
			}
			if nic != nil && utils.StringInSlice(nic.NetworkID,manager.vxnet) {
				manager.cnicache.Enqueue(nic)
				log.Debugf("recycled nic.%s", msg.Nicid)
				manager.qingStub.Tell(qingactor.ModifyNicNameMessage{Nicid: nic.EndpointID, Nicname: manager.cachens})
				ctx.Respond(msg)
				return
			}
		}
		if nic != nil {
			manager.cnicache.RemoveEntries(&msg.Nicid)
			log.Debugf("Removed nic %s from dict",msg.Nicid)
		}
		manager.qingStub.Tell(qingactor.DeleteNicMessage{
				Nic: []*string{&msg.Nicid},
		})
		log.Debugf("cni pool is full,delete nic %s",msg.Nicid)
		ctx.Respond(msg)
	case qingactor.CreateNicReplyMessage:
		manager.policy.UpdateResult(msg.Err)
		unusedList:=manager.cnicache.Enqueue(msg.Nic...)
		if len(unusedList) >0 {
			manager.qingStub.Tell(qingactor.DeleteNicMessage{Nic: unusedList})
			log.Debugf("Deleted unused %d nic",len(unusedList))
		}
	case *actor.Stopping:
		ctx.PopBehavior()
	}
}

func NewNicManager() actor.Actor {
	return &NicManager{}
}

func (manager *NicManager) createOneNewNic(name string, retry int) (response *qingactor.CreateNicReplyMessage, err error) {
	for times := 0; times <= retry; times++ {
		response, err = manager.createNewNicInVxnetSync(name, manager.vxnet[manager.policy.GetNextItem()], 1)
		manager.policy.UpdateResult(err)
		if err == nil {
			return
		}
	}
	return
}

func (manager *NicManager) createNewNicInVxnetSync(nicname string, vxnet string, quantity int) (*qingactor.CreateNicReplyMessage, error) {

	var err error
	var result interface{}
	result, err = manager.qingStub.RequestFuture(qingactor.CreateNicMessage{
		NetworkID: vxnet,
		Nicname:   nicname,
		Quantity:  quantity,
	}, qingactor.DefaultCreateNicTimeout).Result()
	response := result.(qingactor.CreateNicReplyMessage)

	if err != nil || response.Err != nil {
		manager.qingStub.RequestFuture(qingactor.DeleteNicMessage{
			Nic: []*string{&response.Nic[0].EndpointID},
		}, qingactor.DefaultDeleteNicTimeout).Wait()
		err = fmt.Errorf("Failed to create new nic: %v,%v", err, response.Err)
		log.Error(err)
		return nil, err
	}
	return &response, err
}

func (manager *NicManager) createNewNicInVxnetAsync(nicname string, vxnet string, quantity int) {
	manager.qingStub.Tell(qingactor.CreateNicMessage{
		NetworkID: vxnet,
		Nicname:   nicname,
		Quantity:  quantity,
	})
}

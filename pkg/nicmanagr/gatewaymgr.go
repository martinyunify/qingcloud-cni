package nicmanagr

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/qingactor"
	"github.com/vishvananda/netlink"
	"fmt"
	"time"
)

const GatewayManagerActorName = "GatewayManager"

const DefaultGetGatewayTimeout = 60*time.Second
func NewGatewayManager()actor.Actor{
	return &GatewayManager{gateway:make(map[string]*common.Network)}
}
func init() {
	props := actor.FromProducer(NewGatewayManager).WithMailbox(mailbox.Bounded(1000))
	_, err := actor.SpawnNamed(props, GatewayManagerActorName)
	if err != nil {
		log.Error(err)
	} else {
		log.Debugf("Gateway Manager is spawned")
	}
}

type GatewayManager struct {
	gateway  map[string]*common.Network
	qingstub *actor.PID
	nsname string
}

type InitGatewayMessage struct{
	Nsname string
}

func (manager *GatewayManager) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case InitGatewayMessage:
		manager.qingstub = actor.NewLocalPID(qingactor.QingCloudActorName)
		niclist,err:=netlink.LinkList()
		if err != nil {
			log.Error(fmt.Errorf("Failed to get list of nics %v",err))
			return
		}
		var nicliststr []*string
		for _,nic := range niclist{
			macaddr:=nic.Attrs().HardwareAddr.String()
			if macaddr != ""{
				nicliststr= append(nicliststr,&macaddr)
			}
		}
		result,err:=manager.qingstub.RequestFuture(qingactor.DescribeNicMessage{
			Nicid: nicliststr,
		},qingactor.DefaultQueryNicTimeout).Result()
		response := result.(qingactor.DescribeNicReplyMessage)
		if err != nil ||response.Err!= nil {
			log.Error(fmt.Errorf("Failed to send nic query %v, %v",err,response.Err))
			return
		}
		var vxnetlist []*string
		for _,nic := range response.Endpoints{
			if _,ok:=manager.gateway[nic.NetworkID];!ok {
				manager.gateway[nic.NetworkID] = &common.Network{
					Gateway: nic.Address,
				}
				vxnetlist = append(vxnetlist,&nic.NetworkID)
				log.Debugf("loaded nic: %s, %s, %s", nic.NetworkID,nic.EndpointID,nic.Address)
			} else {
				manager.qingstub.RequestFuture(&messages.DeleteNicMessage{Nicid:nic.EndpointID},DefaultDeleteNicTimeout).Wait()
			}
		}

		result, err = manager.qingstub.RequestFuture(qingactor.DescribeVxnetMessage{VxNets: vxnetlist}, qingactor.DefaultQueryVxnetTimeout).Result()
		if err != nil {
			log.Error(fmt.Errorf("Failed to send vxnet query %v, %v", err))
			return
		}
		vxnets:= result.(qingactor.DescribeVxnetReplyMessage)
		if vxnets.Err != nil {
			log.Error(fmt.Errorf("Failed to send vxnet query %v, %v", vxnets.Err))
			return
		}

		for _,vxnet := range vxnets.Networks{
			vxnet.Gateway=manager.gateway[vxnet.NetworkID].Gateway
			manager.gateway[vxnet.NetworkID]=vxnet
		}
		ctx.PushBehavior(manager.ProcessEvent)
		log.Debugf("Gateway manager is initialized")
	}
}

func (manager *GatewayManager) ProcessEvent(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.GetGatewayMessage:
		log.Debugf("Request gateway for vxnet:%s",msg.Vxnet)
		if gateway, ok := manager.gateway[msg.Vxnet]; ok {
			reply := messages.GetGatewayReplyMessage{
				NetworkCIDR:  gateway.Pool,
				Gateway:    gateway.Gateway,
				NetworkID: gateway.NetworkID,
			}
			context.Respond(&reply)
		} else {
			result,err:=manager.qingstub.RequestFuture(qingactor.CreateNicMessage{
				NetworkID:msg.Vxnet,
				Nicname:"hostnic",
			},qingactor.DefaultCreateNicTimeout).Result()
			if err != nil {
				context.Respond(&messages.GetGatewayReplyMessage{})
				log.Error(fmt.Errorf("Failed to create new nic %v",err))
				return
			}
			response := result.(qingactor.CreateNicReplyMessage)
			if response.Err != nil {
				context.Respond(&messages.GetGatewayReplyMessage{})
				log.Error(fmt.Errorf("Failed to create new nic %v",response.Err))
			}
			result,err= manager.qingstub.RequestFuture(qingactor.DescribeVxnetMessage{VxNets:[]*string{&response.Nic.NetworkID}},qingactor.DefaultQueryVxnetTimeout).Result()
			if err != nil {
				context.Respond(&messages.GetGatewayReplyMessage{})
				log.Error(fmt.Errorf("Failed to get vxnet %v",err))
				return
			}
			vxnetresp:= result.(qingactor.DescribeVxnetReplyMessage)
			if vxnetresp.Err != nil {
				context.Respond(&messages.GetGatewayReplyMessage{})
				log.Error(fmt.Errorf("Failed to get vxnet %v",vxnetresp.Err))
			}
			vxnetresp.Networks[0].Gateway = response.Nic.Address
			manager.gateway[response.Nic.NetworkID]= vxnetresp.Networks[0]
			context.Respond(&messages.GetGatewayReplyMessage{
				NetworkID:response.Nic.NetworkID,
				NetworkCIDR:vxnetresp.Networks[0].Pool,
				Gateway: response.Nic.Address,
			})
		}
	}
}


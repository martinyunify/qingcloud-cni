package nicmanagr

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-cni/pkg/qingactor"
	"time"
	"net"
	"github.com/yunify/qingcloud-cni/pkg/utils"
)

const GatewayManagerActorName = "GatewayManager"

const DefaultGetGatewayTimeout = 60 * time.Second

func NewGatewayManager() actor.Actor {
	return &GatewayManager{gateway: make(map[string]*common.Network)}
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
	vxnet []string
	nsname   string
}

type InitGatewayMessage struct {
	Nsname string
	Vxnet []string
}

func (manager *GatewayManager) Receive(ctx actor.Context) {
	switch msg:= ctx.Message().(type) {
	case InitGatewayMessage:
		manager.vxnet = msg.Vxnet
		manager.nsname = msg.Nsname
		manager.qingstub = actor.NewLocalPID(qingactor.QingCloudActorName)
		niclist, err := netlink.LinkList()
		if err != nil {
			log.Error(fmt.Errorf("Failed to get list of nics %v", err))
			return
		}
		var nicliststr []*string

		for _, nic := range niclist {
			macaddr := nic.Attrs().HardwareAddr.String()
			if macaddr != "" {
				nicliststr = append(nicliststr, &macaddr)
			}
		}
		result, err := manager.qingstub.RequestFuture(qingactor.DescribeNicMessage{
			Nicid: nicliststr,
		}, qingactor.DefaultQueryNicTimeout).Result()
		response := result.(qingactor.DescribeNicReplyMessage)
		if err != nil || response.Err != nil {
			log.Error(fmt.Errorf("Failed to send nic query %v, %v", err, response.Err))
			return
		}


		var unusedList []*string

		//inefficient way
		for _, nic := range response.Endpoints {
			if utils.StringInSlice(nic.NetworkID,manager.vxnet) {
				if _, ok := manager.gateway[nic.NetworkID]; !ok {
					niclink, err := utils.LinkByMacAddr(nic.EndpointID)
					if err == nil&&niclink.Attrs().Flags|net.FlagUp != 0 {
						manager.gateway[nic.NetworkID] = &common.Network{
							Gateway: nic.Address,
						}
						log.Debugf("loaded nic: %s, %s, %s", nic.NetworkID, nic.EndpointID, nic.Address)
						continue
					}
				}
			}
			unusedList = append(unusedList, &nic.EndpointID)
		}

		var vxnetList []*string
		for _,vnet := range manager.vxnet{
			vxnetList = append(vxnetList,&vnet)
		}
		result, err = manager.qingstub.RequestFuture(qingactor.DescribeVxnetMessage{VxNets: vxnetList}, qingactor.DefaultQueryVxnetTimeout).Result()
		if err != nil {
			log.Error(fmt.Errorf("Failed to send vxnet query %v, %v", err))
			return
		}
		vxnets := result.(qingactor.DescribeVxnetReplyMessage)
		if vxnets.Err != nil {
			log.Error(fmt.Errorf("Failed to send vxnet query %v, %v", vxnets.Err))
			return
		}

		for _, vxnet := range vxnets.Networks {
			if _,ok:=manager.gateway[vxnet.NetworkID]; ok {
				vxnet.Gateway = manager.gateway[vxnet.NetworkID].Gateway
			} else {
				manager.createNewNicInVxnetAsync(manager.nsname,vxnet.NetworkID,1)
			}
			manager.gateway[vxnet.NetworkID] = vxnet
		}
		manager.qingstub.Tell(qingactor.DeleteNicMessage{Nic:unusedList})

	case qingactor.CreateNicReplyMessage:
		for _,nic := range msg.Nic {
			manager.gateway[nic.NetworkID].Gateway = nic.Address
		}

		if len(manager.gateway) == len(manager.vxnet) {

			ctx.PushBehavior(manager.ProcessEvent)
			log.Debugf("Gateway manager is initialized")
		} else {
			for _,vxnet := range manager.vxnet {
				if _,ok:=manager.gateway[vxnet]; !ok {
					manager.createNewNicInVxnetAsync(manager.nsname,vxnet,1)
				}
			}
		}
	}
}

func (manager *GatewayManager) ProcessEvent(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.GetGatewayMessage:
		log.Debugf("Request gateway for vxnet:%s", msg.Vxnet)
		if gateway, ok := manager.gateway[msg.Vxnet]; ok {
			reply := messages.GetGatewayReplyMessage{
				NetworkCIDR: gateway.Pool,
				Gateway:     gateway.Gateway,
				NetworkID:   gateway.NetworkID,
			}
			context.Respond(&reply)
		} else {
			log.Errorf("Vxnet %s is not in vxnet list %v", msg.Vxnet, manager.vxnet)
			context.Respond(&messages.GetGatewayReplyMessage{})
		}
	}
}

func (manager *GatewayManager) createNewNicInVxnetSync(nicname string, vxnet string, quantity int, retry int) (response qingactor.CreateNicReplyMessage,err error) {
	for times := 0; times <= retry; times++ {
		var result interface{}
		result, err = manager.qingstub.RequestFuture(qingactor.CreateNicMessage{
			NetworkID: vxnet,
			Nicname:   nicname,
			Quantity:  quantity,
		}, qingactor.DefaultCreateNicTimeout).Result()
		response = result.(qingactor.CreateNicReplyMessage)

		if err != nil || response.Err != nil {
			manager.qingstub.RequestFuture(qingactor.DeleteNicMessage{
				Nic: []*string{&response.Nic[0].EndpointID},
			}, qingactor.DefaultDeleteNicTimeout).Wait()
			err = fmt.Errorf("Failed to create new nic: %v,%v", err, response.Err)
			log.Error(err)
			continue
		}
		return response,err
	}
	return
}

func (manager *GatewayManager) createNewNicInVxnetAsync(nicname string, vxnet string, quantity int) {
	manager.qingstub.Tell(qingactor.CreateNicMessage{
		NetworkID: vxnet,
		Nicname:   nicname,
		Quantity:  quantity,
	})
}
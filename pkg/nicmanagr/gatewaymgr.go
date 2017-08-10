package nicmanagr

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-cni/pkg/messages"
)

type GatewayManager struct {
	gateway  map[string]common.Endpoint
	iface    common.Endpoint
	qingstub *actor.PID
}

func (manager *GatewayManager) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.QingcloudInitializeMessage:

		manager.loadDefaultGateway(msg.Iface)
	}
}

func (manager *GatewayManager) ProcessEvent(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.GetGatewayMessage:
		if gateway, ok := manager.gateway[msg.Vxnet]; ok {
			reply := messages.GetGatewayReplyMessage{
				NetworkID:  gateway.NetworkID,
				EndpointID: gateway.EndpointID,
				Address:    gateway.Address,
			}
			context.Respond(reply)
		}
	}
}
func (manager *GatewayManager) loadDefaultGateway(iface string) {

}

package qingactor

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-sdk-go/service"
)

//RouterActor Router actor for qingcloud
type RouterActor struct {
	routerStub *service.RouterService
	jobStub *service.JobService

}

//JoinRouterMessage join Router and allocate vxnet
type JoinRouterMessage struct {
	network common.Network
	router  common.Router
}

//JoinRouterReplyMessage join Router and allocate vxnet reply
type JoinRouterReplyMessage struct {
	network common.Network
	router  common.Router
}

type LeaveRouterMessage struct {
	network common.Network
	router  common.Router
}

func (qactor *RouterActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case JoinRouterMessage:
		gatewayID := msg.network.Ipv4Data.Gateway.String()
		request := service.JoinRouterInput{
			VxNet:     &msg.network.NetworkID,
			Router:    &msg.router.RouterID,
			IPNetwork: &msg.network.Ipv4Data.AddressSpace,
			ManagerIP: &gatewayID,
		}
		qactor.routerStub.JoinRouter(&request)

	case LeaveRouterMessage:

	}
}

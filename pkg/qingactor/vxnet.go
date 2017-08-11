package qingactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-sdk-go/service"
)

//VxNetActor vxnet qingcloud handler
type VxNetActor struct {
	vxNetStub *service.VxNetService
	jobStub   *service.JobService
}

//CreateVxNetMessage create vxnet, optional parameter NetworkID: identify network using custom id
type CreateVxNetMessage struct {
	NetworkID string
}

//CreateVxNetReplyMessage create vxnet reply
type CreateVxNetReplyMessage struct {
	err     error
	network common.Network
}

//DeleteVxnetMessage delete vxnet input
type DeleteVxnetMessage struct {
	network common.Network
}

//DeleteVxnetReplyMessage delete vxnet output
type DeleteVxnetReplyMessage struct {
	err error
}

type DescribeVxnetMessage struct {
	VxNets []*string
}

type DescribeVxnetReplyMessage struct {
	Err error
}

//Receive message handler function
func (vxnet *VxNetActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case CreateVxNetMessage:
		counter := 1
		vxnetType := 0
		input := service.CreateVxNetsInput{
			Count:     &counter,
			VxNetName: &msg.NetworkID,
			VxNetType: &vxnetType,
		}
		result, err := vxnet.vxNetStub.CreateVxNets(&input)
		reply := CreateVxNetReplyMessage{}
		if err != nil || *result.RetCode != 0 {
			reply.err = err
			if err == nil {
				reply.err = fmt.Errorf("Failed to create vxnet %s", *result.Message)
			}
			log.Errorln(reply.err)
		} else {
			vxnet := result.VxNets[0]
			reply.network = common.Network{
				NetworkID:   msg.NetworkID,
				NetworkName: *vxnet,
			}
		}
		context.Respond(reply)
	case DeleteVxnetMessage:
		request := service.DeleteVxNetsInput{
			VxNets: []*string{&msg.network.NetworkName},
		}
		result, err := vxnet.vxNetStub.DeleteVxNets(&request)
		reply := DeleteVxnetReplyMessage{}
		if err != nil || *result.RetCode != 0 {
			reply.err = err
			if err == nil {
				reply.err = fmt.Errorf("Failed to delete vxnet: %s", *result.Message)
			}
			log.Errorln(reply.err)
		} else {
			reply.err = nil
		}
		context.Respond(reply)
	case DescribeVxnetMessage:
		request := service.DescribeVxNetsInput{
			VxNets: msg.VxNets,
		}

		result, err := vxnet.vxNetStub.DescribeVxNets(&request)
		reply := DescribeVxnetReplyMessage{}
		if err != nil || *result.RetCode != 0 {
			reply.Err = err
			if err == nil {
				reply.Err = fmt.Errorf("Failed to describe vxnet: %s", *result.Message)
			}
		} else {
			reply.Err = nil
		}
		context.Respond(reply)
	}
}

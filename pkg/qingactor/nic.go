package qingactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/container-nanny/pkg/common"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-sdk-go/service"
)

//NicActor nic qingcloud handler
type NicActor struct {
	nicStub *service.NicService
	zone    string
}

//CreateNicMessage create nic message
type CreateNicMessage struct {
	NetworkID  string
	EndpointID string
	Address    string
}

//CreateNicReplyMessage create vxnet reply message
type CreateNicReplyMessage struct {
	err error
	nic common.Endpoint
}

type DeleteNicMessage struct {
	nic common.Endpoint
}

type DeleteNicReplyMessage struct {
	err error
}

type DescribeNicMessage struct {
	instanceid string
	nicid      string
}

type DescirbeNicReplyMessage struct {
}

//Receive message handler function
func (nic *NicActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.AllocateNicMessage:
		count := 1
		request := service.CreateNicsInput{
			Count:   &count,
			NICName: &msg.Name,
			VxNet:   &msg.NetworkID,
		}
		result, err := nic.nicStub.CreateNics(&request)
		reply := CreateNicReplyMessage{}
		if err != nil || *result.RetCode != 0 {
			reply.err = err
			if err == nil {
				reply.err = fmt.Errorf("Failed to create nic:%s", *result.Message)
			}
			log.Error(reply.err)
		} else {
			nic := result.Nics[0]
			reply.nic = common.Endpoint{
				NetworkID:  msg.NetworkID,
				EndpointID: *nic.NICID,
			}
		}
		context.Respond(reply)
	case DeleteNicMessage:
		request := service.DeleteNicsInput{
			Nics: []*string{&msg.nic.EndpointID},
		}
		result, err := nic.nicStub.DeleteNics(&request)
		if err == nil {
			err = fmt.Errorf("Failed to delete nic: %s", result.Message)
		}
		reply := DeleteNicReplyMessage{
			err: err,
		}
		context.Respond(reply)
	case DescribeNicMessage:
		request := service.DescribeNicsInput{}
	}
}

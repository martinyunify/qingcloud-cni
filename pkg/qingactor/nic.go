package qingactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-sdk-go/service"
)

//NicActor nic qingcloud handler
type NicActor struct {
	nicStub *service.NicService
	zone    string
}

//CreateNicMessage create nic message
type CreateNicMessage struct {
	NetworkID string
	Nicname   string
}

//CreateNicReplyMessage create vxnet reply message
type CreateNicReplyMessage struct {
	Err error
	Nic common.Endpoint
}

type DeleteNicMessage struct {
	Nic string
}

type DeleteNicReplyMessage struct {
	Err error
}

type DescribeNicMessage struct {
	instanceid string
	nicid      string
}

type DescirbeNicReplyMessage struct {
}

//Receive message handler function
func (nicactor *NicActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case CreateNicMessage:
		count := 1
		request := service.CreateNicsInput{
			Count:   &count,
			NICName: &msg.Nicname,
			VxNet:   &msg.NetworkID,
		}
		result, err := nicactor.nicStub.CreateNics(&request)
		reply := CreateNicReplyMessage{}
		if err != nil || *result.RetCode != 0 {
			reply.Err = err
			if err == nil {
				reply.Err = fmt.Errorf("Failed to create nic:%s", *result.Message)
			}
			log.Error(reply.Err)
			context.Respond(reply)
			return
		} else {
			nic := result.Nics[0]
			reply.Nic = common.Endpoint{
				NetworkID:  msg.NetworkID,
				EndpointID: *nic.NICID,
				Address:    *nic.PrivateIP,
			}

			instanceid, err := loadInstanceID()
			if err != nil {
				reply.Err = fmt.Errorf("Failed to load instanceid: %v", err)
				log.Error(reply.Err)
				context.Respond(reply)
				return
			}
			request := service.AttachNicsInput{
				Nics:     []*string{nic.NICID},
				Instance: &instanceid,
			}
			result, err := nicactor.nicStub.AttachNics(&request)
			if err != nil || *result.RetCode != 0 {
				reply.Err = err
				context.Respond(reply)
				return
			}
			context.Respond(reply)

		}
	case DeleteNicMessage:
		request := service.DetachNicsInput{
			Nics: []*string{&msg.Nic},
		}
		reply := DeleteNicReplyMessage{}

		result, err := nicactor.nicStub.DetachNics(&request)
		if err != nil || *result.RetCode != 0 {
			reply.Err = err
			if reply.Err == nil {
				reply.Err = fmt.Errorf("Failed to detach nic:%s", *result.Message)
			}
			context.Respond(reply)
			return
		} else {
			request := service.DeleteNicsInput{
				Nics: []*string{&msg.Nic},
			}
			result, err := nicactor.nicStub.DeleteNics(&request)
			if err != nil || *result.RetCode != 0 {
				reply.Err = err
				if reply.Err == nil {
					err = fmt.Errorf("Failed to delete nic: %s", result.Message)
				}
			}
			context.Respond(reply)
		}
	}
}

package qingactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/common"
	"github.com/yunify/qingcloud-sdk-go/client"
	"github.com/yunify/qingcloud-sdk-go/service"

	"github.com/yunify/qingcloud-cni/pkg/utils"
	qcutil "github.com/yunify/qingcloud-sdk-go/utils"
	"net"
	"time"
)

//NicActor nic qingcloud handler
type NicActor struct {
	nicStub *service.NicService
	jobStub *service.JobService
}

//CreateNicMessage create nic message
type CreateNicMessage struct {
	NetworkID string
}

//CreateNicReplyMessage create vxnet reply message
type CreateNicReplyMessage struct {
	Err error
	Nic common.Endpoint
}

type DeleteNicMessage struct {
	Nic     string
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

const (
	DefaultCreateNicTimeout = 30 * time.Second
	DefaultDeleteNicTimeout = 30 * time.Second
)

//Receive message handler function
func (nicactor *NicActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case CreateNicMessage:
		count := 1
		request := service.CreateNicsInput{
			Count:   &count,
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
		}
		nic := result.Nics[0]
		reply.Nic = common.Endpoint{
			NetworkID:  msg.NetworkID,
			EndpointID: *nic.NICID,
			Address:    *nic.PrivateIP,
		}
		log.Debugf("Created nic %s", *nic.NICID)

		//attach nic to host
		instanceid, err := loadInstanceID()
		if err != nil {
			reply.Err = fmt.Errorf("Failed to load instanceid: %v", err)
			log.Error(reply.Err)
			context.Respond(reply)
			return
		}
		attrequest := service.AttachNicsInput{
			Nics:     []*string{nic.NICID},
			Instance: &instanceid,
		}
		attresult, err := nicactor.nicStub.AttachNics(&attrequest)
		if err != nil || *attresult.RetCode != 0 {
			reply.Err = err
			if reply.Err == nil {
				reply.Err = fmt.Errorf("failed to attach nic %s", attresult.Message)
			}
		}

		if err = nicactor.waitNic(*attresult.JobID, *nic.NICID, net.FlagUp); err != nil {
			reply.Err = err
		}

		context.Respond(reply)
		log.Debugf("Attached nic %s to host", *nic.NICID)

	case DeleteNicMessage:
		reply := DeleteNicReplyMessage{}

		request := service.DetachNicsInput{
			Nics: []*string{&msg.Nic},
		}

		detresult, err := nicactor.nicStub.DetachNics(&request)
		if err != nil || *detresult.RetCode != 0 {
			reply.Err = err
			if reply.Err == nil {
				reply.Err = fmt.Errorf("Failed to detach nic %s", detresult.Message)
			}
			context.Respond(reply)
			return
		}
		if err = client.WaitJob(nicactor.jobStub, *detresult.JobID, 25*time.Second, 5*time.Second); err != nil {
			reply.Err = err
			context.Respond(reply)
			return
		}
		delRequest := service.DeleteNicsInput{
			Nics: []*string{&msg.Nic},
		}

		delresult, err := nicactor.nicStub.DeleteNics(&delRequest)
		if err != nil || *delresult.RetCode != 0 {
			reply.Err = err
			if reply.Err == nil {
				reply.Err = fmt.Errorf("Failed to delete nic: %s", delresult.Message)
			}
		}
		context.Respond(reply)
	}
}

func (nicactor *NicActor) waitNic(jobid string, nicID string, status net.Flags) error {
	log.Debugf("Wait for nic %v", nicID)
	err := qcutil.WaitForSpecific(func() bool {
		link, err := utils.LinkByMacAddr(nicID)
		if err != nil {
			return false
		}
		if link.Attrs().Flags|status != 0 {
			log.Debugf("Find link %s %s", link.Attrs().Name, nicID)
			return true
		}
		return false
	}, 25*time.Second, 5*time.Second)
	if err != nil {
		return err
	}
	return client.WaitJob(nicactor.jobStub, jobid, 25*time.Second, 1*time.Second)
}

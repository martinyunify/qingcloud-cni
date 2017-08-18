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
	Nicname   string
	Quantity  int
}

//CreateNicReplyMessage create vxnet reply message
type CreateNicReplyMessage struct {
	Err error
	Nic []*common.Endpoint
}

type DeleteNicMessage struct {
	Nic     []*string
}

type DeleteNicReplyMessage struct {
	Err error
}

type DescribeNicMessage struct {
	Instanceid string
	Nicid      []*string
	Nicname    string
}

type DescribeNicReplyMessage struct {
	Err       error
	Endpoints []*common.Endpoint
}

type ModifyNicNameMessage struct{
	Nicid string
	Nicname string
}

const (
	DefaultCreateNicTimeout = 30 * time.Second
	DefaultDeleteNicTimeout = 30 * time.Second
	DefaultQueryNicTimeout  = 30 * time.Second
)

//Receive message handler function
func (nicactor *NicActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case CreateNicMessage:
		request := service.CreateNicsInput{
			Count:   &msg.Quantity,
			VxNet:   &msg.NetworkID,
			NICName: &msg.Nicname,
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
		for _, nic := range result.Nics {
			reply.Nic = append(reply.Nic, &common.Endpoint{
				NetworkID:  msg.NetworkID,
				EndpointID: *nic.NICID,
				Address:    *nic.PrivateIP,
			})
			log.Debugf("Created nic %s", *nic.NICID)
		}

		//attach nic to host
		instanceid, err := loadInstanceID()
		if err != nil {
			reply.Err = fmt.Errorf("Failed to load instanceid: %v", err)
			log.Error(reply.Err)
			context.Respond(reply)
			return
		}
		var nicidlist []*string
		for _, nic := range result.Nics {
			nicidlist = append(nicidlist, nic.NICID)
		}
		attrequest := service.AttachNicsInput{
			Nics:     nicidlist,
			Instance: &instanceid,
		}
		attresult, err := nicactor.nicStub.AttachNics(&attrequest)
		if err != nil || *attresult.RetCode != 0 {
			reply.Err = err
			if reply.Err == nil {
				reply.Err = fmt.Errorf("failed to attach nic %s", attresult.Message)
			}
		}

		if err = nicactor.waitNic(*attresult.JobID, nicidlist, net.FlagUp); err != nil {
			reply.Err = err
		}

		context.Respond(reply)
		log.Debugf("Attached nic %s to host", nicidlist)

	case DeleteNicMessage:
		reply := DeleteNicReplyMessage{}

		request := service.DetachNicsInput{
			Nics: msg.Nic,
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
			Nics: msg.Nic,
		}

		delresult, err := nicactor.nicStub.DeleteNics(&delRequest)
		if err != nil || *delresult.RetCode != 0 {
			reply.Err = err
			if reply.Err == nil {
				reply.Err = fmt.Errorf("Failed to delete nic: %s", delresult.Message)
			}
		}
		context.Respond(reply)
	case DescribeNicMessage:
		reply := DescribeNicReplyMessage{}

		request := service.DescribeNicsInput{}
		if msg.Nicname != "" {
			request.NICName = &msg.Nicname
		}
		if len(msg.Nicid) > 0 {
			request.Nics = msg.Nicid
		}
		if msg.Instanceid != "" {
			request.Instances = []*string{&msg.Instanceid}
		}
		descriResult, err := nicactor.nicStub.DescribeNics(&request)
		if err != nil {
			reply.Err = fmt.Errorf("Failed to describe nic:%v", err)
			context.Respond(reply)
			return
		}
		if *descriResult.RetCode != 0 {
			reply.Err = fmt.Errorf("Failed to describe nic: %s", *descriResult.Message)
			context.Respond(reply)
			return
		}
		for _, nic := range descriResult.NICSet {
			reply.Endpoints = append(reply.Endpoints, &common.Endpoint{
				Address:    *nic.PrivateIP,
				EndpointID: *nic.NICID,
				NetworkID:  *nic.VxNetID,
			})
		}
		context.Respond(reply)
	case ModifyNicNameMessage:
		request:= service.ModifyNicAttributesInput{
			NICID:&msg.Nicid,
			NICName:&msg.Nicname,
		}
		nicactor.nicStub.ModifyNicAttributes(&request)
	}
}

func (nicactor *NicActor) waitNic(jobid string, nicIDs []*string, status net.Flags) error {
	log.Debugf("Wait for nic %v", nicIDs)
	err := qcutil.WaitForSpecific(func() bool {
		for _, nicID := range nicIDs {
			link, err := utils.LinkByMacAddr(*nicID)
			if err != nil {
				return false
			}
			if link.Attrs().Flags|status == 0 {
				return false
			}
			log.Debugf("Find link %s %s", link.Attrs().Name, *nicID)
		}
		return true
	}, 25*time.Second, 1*time.Second)
	if err != nil {
		return err
	}
	return client.WaitJob(nicactor.jobStub, jobid, 25*time.Second, 1*time.Second)
}

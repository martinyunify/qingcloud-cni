package qingactor

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"
	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/service"
	"io/ioutil"
)

const (
	instanceIDFile = "/etc/qingcloud/instance-id"
)

//QingCloudActor QingCloudService Actor
type QingCloudActor struct {
	routerStub *actor.PID
	vxNetStub  *actor.PID
	nicStub    *actor.PID
	zone       string
}

func NewQingCloudActor() actor.Actor {
	return &QingCloudActor{}
}

//Stop stop all of actors
func (qactor *QingCloudActor) Stop() {
	if qactor.routerStub != nil {
		qactor.routerStub.GracefulStop()
	}
	if qactor.vxNetStub != nil {
		qactor.vxNetStub.GracefulStop()
	}
	if qactor.nicStub != nil {
		qactor.nicStub.GracefulStop()
	}
}

func (qactor *QingCloudActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case messages.QingcloudInitializeMessage:
		var err error
		newConfig, err := config.NewDefault()
		if err != nil {
			log.Errorf("Default qingcloud newConfig is not created")
			return
		}
		newConfig.LoadConfigFromFilepath(msg.QingAccessFile)
		newConfig.Zone = msg.Zone
		qactor.zone = msg.Zone
		qingStub, err := service.Init(newConfig)
		if err != nil {
			log.Errorf("Failed to read qingcloud newConfig file")
			return
		}

		routerStub, err := qingStub.Router(qactor.zone)
		if err != nil {
			log.Errorf("Failed to create router stub")
			return
		}
		props := actor.FromInstance(&RouterActor{routerStub: routerStub, zone: qactor.zone})
		qactor.routerStub = actor.Spawn(props)

		vxnetStub, err := qingStub.VxNet(qactor.zone)
		if err != nil {
			log.Errorf("Failed to create vxnet stub")
			qactor.Stop()
			return
		}
		props = actor.FromInstance(&VxNetActor{vxNetStub: vxnetStub, zone: qactor.zone})
		qactor.vxNetStub = actor.Spawn(props)

		nicStub, err := qingStub.Nic(qactor.zone)
		if err != nil {
			log.Errorf("Failed to create nic stub")
			qactor.Stop()
			return
		}

		props = actor.FromInstance(&NicActor{nicStub: nicStub, zone: qactor.zone})
		qactor.nicStub = actor.Spawn(props)
		log.Debugf("QingCloud sdk is initialized.")

		context.PushBehavior(qactor.ProcessMsg)
	}
}

//Receive REceive actor messages
func (qactor *QingCloudActor) ProcessMsg(context actor.Context) {
	switch msg := context.Message().(type) {
	case CreateVxNetMessage:
		qactor.vxNetStub.Request(msg, context.Sender())
	case CreateNicMessage:
		qactor.nicStub.Request(msg, context.Sender())
	case DeleteNicMessage:
		qactor.nicStub.Request(msg, context.Sender())
	case *actor.Stopping:
		qactor.Stop()
	}
}

func loadInstanceID() (string, error) {
	content, err := ioutil.ReadFile(instanceIDFile)
	if err != nil {
		return "", fmt.Errorf("Load instance-id from %s error: %v", instanceIDFile, err)
	}
	return string(content), nil
}

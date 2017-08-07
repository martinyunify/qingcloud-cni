package qingcloud

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/yunify/qingcloud-cni/pkg/messages"
)

type MockQingCloud struct {
}

func (stub *MockQingCloud) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case messages.AllocateNicMessage:
		ctx.Respond(messages.AllocateNicReplyMessage{Name: msg.GetName()})
	}
}

func NewMockQingCloud() actor.Actor {
	return &MockQingCloud{}
}

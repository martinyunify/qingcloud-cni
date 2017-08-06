package qingcloud

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/yunify/qingcloud-cni/pkg/messages"
)

type MockQingCloud struct {
}

func (stub *MockQingCloud) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case messages.CreateNewNicMessage:
		ctx.Respond(messages.ErrorMessage{Err: fmt.Errorf("mock Replied,%v", msg)})
	}
}

func NewMockQingCloud() actor.Actor {
	return &MockQingCloud{}
}

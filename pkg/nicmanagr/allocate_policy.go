package nicmanagr

import (
	"fmt"
	"time"
	"math/rand"
	"strings"
)

type CreationPolicy struct {
	poolSize int
	policyFuction PolicyFunction
	flag int
	lastError error
}

type PolicyFunction func(size int,flag int,laststatus error) int

type PolicyType uint
const (
	RoundRotate = iota
	FailRotate
	Random
	Unsupported
)

func GetPolicyType(policy string)PolicyType{
	switch strings.ToLower(policy) {
	case "roundrotate":
		return RoundRotate
	case "failrotate":
		return FailRotate
	case "ramdom":
		return Random
	default:
		return Unsupported
	}
}

func NewCreationPolicy(size int, policytype PolicyType) (policy *CreationPolicy,err error){
	policy = &CreationPolicy{
		flag:0,
	}
	if err = policy.SetPoolSize(size); err != nil{
		return nil, err
	}
	if err = policy.SetPolicy(policytype); err != nil{
		return nil,err
	}
	return
}

func (policy *CreationPolicy) SetPoolSize(size int) error{
	if size <=0 {
		return fmt.Errorf("Resource pool size should be greater than 0")
	}
	policy.poolSize = size
	return nil
}

func (policy *CreationPolicy) SetPolicy(policytype PolicyType) error{
	switch policytype {
	case RoundRotate:
		policy.policyFuction = policy.roundRotate
	case FailRotate:
		policy.policyFuction = policy.failRotate
	case Random:
		policy.policyFuction = policy.random
	default:
		return fmt.Errorf("Unsupported Policy")
	}
	return nil
}

func (policy *CreationPolicy) UpdateResult(err error){
	policy.lastError = err
}

func (policy *CreationPolicy) GetNextItem()int {
	if policy.policyFuction != nil {
		policy.flag=policy.policyFuction(policy.poolSize,policy.flag,policy.lastError)
	}
	return policy.flag
}

func (policy *CreationPolicy) failRotate(size int,flag int,laststatus error)int {
	if laststatus == nil {
		return flag
	}
	return flag+1 % size
}

func (policy *CreationPolicy) roundRotate(size int,flag int,laststatus error)int {
	return flag+1 % size
}

func (policy *CreationPolicy) random(size int,flag int,laststatus error)int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(size)
}

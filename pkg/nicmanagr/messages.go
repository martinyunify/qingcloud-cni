package nicmanagr

import (
	"fmt"
)


type ResourcePoolInitMessage struct{
	Vxnet []string
	Policy PolicyType
}

func NewResourcePoolInitMessage(vxnet []string,policy string)(*ResourcePoolInitMessage,error){
	for _, item := range vxnet {
		if item == "" || item == "vxnet-xxxxxxx" {
			return nil, fmt.Errorf("Invalid vxnet: %s", item)
		}
	}

	policytype :=GetPolicyType(policy)
	if policytype == Unsupported {
		return nil,fmt.Errorf("Invalid policy")
	}
	return &ResourcePoolInitMessage{
		Vxnet:vxnet,
		Policy: policytype,
	},nil
}

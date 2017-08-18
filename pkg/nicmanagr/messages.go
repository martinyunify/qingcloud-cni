package nicmanagr

import (
	"fmt"
)

type ResourcePoolInitMessage struct {
	Vxnet          []string
	Policy         PolicyType
	NicNameinCache string
	NicCacheSize   int
}

func NewResourcePoolInitMessage(vxnet []string, policy string, NicNameinCache string, nicCachesize int) (*ResourcePoolInitMessage, error) {
	for _, item := range vxnet {
		if item == "" || item == "vxnet-xxxxxxx" {
			return nil, fmt.Errorf("Invalid vxnet: %s", item)
		}
	}

	policytype := GetPolicyType(policy)
	if policytype == Unsupported {
		return nil, fmt.Errorf("Invalid policy")
	}

	if nicCachesize < 0 {
		return nil, fmt.Errorf("Cache Size should be greater than 0")
	}

	return &ResourcePoolInitMessage{
		Vxnet:          vxnet,
		Policy:         policytype,
		NicNameinCache: NicNameinCache,
		NicCacheSize: nicCachesize,
	}, nil
}

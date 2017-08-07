package nicmanagr

import "github.com/yunify/qingcloud-cni/pkg/common"

type gatewaymgr struct {
	gateway map[string]common.Endpoint
	iface   string
}

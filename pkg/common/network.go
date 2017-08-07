package common

import (
	"net"
)

type Network struct {
	NetworkID   string
	NetworkName string
	Options     map[string]interface{}
	Ipv4Data    *IPAMData
}

type IPAMData struct {
	Pool         *net.IPNet
	Gateway      *net.IPAddr
	AddressSpace string
	AuxAddresses map[string]*net.IPAddr
}

type VPC struct {
	VPCID    string
	networks []*Network
}

type Endpoint struct {
	NetworkID  string
	EndpointID string
	Interface  *EndpointInterface
	Options    map[string]interface{}
}

type EndpointInterface struct {
	Address     *net.IPAddr
	AddressIPv6 *net.IPAddr
	MacAddress  *net.HardwareAddr
}

//Router router struct
type Router struct {
	RouterID string
}

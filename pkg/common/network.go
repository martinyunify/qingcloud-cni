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
	Address    string
}

//Router router struct
type Router struct {
	RouterID string
}

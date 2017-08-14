package common

type Network struct {
	NetworkID   string
	NetworkName string
	Pool        string
	Gateway     string
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

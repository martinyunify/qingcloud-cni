package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
)

func cmdDel(args *skel.CmdArgs) error {
	return nil
}

func cmdAdd(args *skel.CmdArgs) error {
	return nil
}

func main(){
	skel.PluginMain(cmdAdd, cmdDel, version.All)

}
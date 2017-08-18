package utils

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"syscall"
)

func LinkByMacAddr(macAddr string) (netlink.Link, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		attr := link.Attrs()
		if attr.HardwareAddr.String() == macAddr {
			return link, nil
		}
	}
	return nil, fmt.Errorf("Can not find link by address: %s", macAddr)
}

func ClearNic(iface netlink.Link)error{
	if iface != nil {
		//clear old addr. os possible bind ip to nic before hostnic.
		addrs, err := netlink.AddrList(iface, syscall.AF_INET)
		if err == nil && len(addrs) > 0 {
			for _, addr := range addrs {
				err := netlink.AddrDel(iface, &addr)
				if err != nil {
					return fmt.Errorf("AddrDel err %s addr:%+v, Nic %s", err.Error(), addr, iface.Attrs().HardwareAddr)
				}
			}
		}
	}
	return nil
}
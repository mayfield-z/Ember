package utils

import (
	"net"
	"os/exec"
)

func Add1(ip net.IP) net.IP {
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] < 255 {
			ip[i]++
			return ip
		}
	}
	return ip
}

func AddIpToInterface(ip net.IP, iface *net.Interface) error {
	if iface == nil {
		return nil
	}
	existAddr, err := iface.Addrs()
	if err != nil {
		return err
	}
	for _, addr := range existAddr {
		if ip.Equal(addr.(*net.IPNet).IP) {
			return nil
		}
	}
	cmd := exec.Command("ip", "addr", "add", ip.String(), "dev", iface.Name)
	return cmd.Run()
}

func DelIpFromInterface(ip net.IP, iface *net.Interface) error {
	if iface == nil {
		return nil
	}
	existAddr, err := iface.Addrs()
	if err != nil {
		return err
	}
	for _, addr := range existAddr {
		if ip.Equal(addr.(*net.IPNet).IP) {
			cmd := exec.Command("ip", "addr", "del", ip.String(), "dev", iface.Name)
			return cmd.Run()
		}
	}
	return nil
}

func DelAddrFromInterface(addr net.Addr, iface *net.Interface) error {
	if iface == nil {
		return nil
	}
	existAddrs, err := iface.Addrs()
	if err != nil {
		return err
	}
	for _, existAddr := range existAddrs {
		if existAddr.String() == addr.String() {
			cmd := exec.Command("ip", "addr", "del", addr.String(), "dev", iface.Name)
			return cmd.Run()
		}
	}
	return nil
}

func ExistsInAddrList(addr net.Addr, addrList []net.Addr) bool {
	for _, addrInList := range addrList {
		if addr.String() == addrInList.String() {
			return true
		}
	}
	return false
}

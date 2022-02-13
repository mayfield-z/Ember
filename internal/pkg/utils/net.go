package utils

import (
	"net"
	"os/exec"
)

func Add1(ip net.IP) net.IP {
	newIp := make([]byte, 16)
	copy(newIp, ip)
	for i := len(newIp) - 1; i >= 0; i-- {
		if newIp[i] < 255 {
			newIp[i]++
			return newIp
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

package tool

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sort"
)

// 获取管理地址
func HostAdmIp(intfList []string) (string, error) {
	if intfList == nil {
		intfList = []string{"bond0", "eth0", "eth1"}
	}
	return hostAdmIp(intfList)
}

func hostAdmIp(intfList []string) (string, error) {
	var (
		intf *net.Interface
		err  error
	)
	for i := 0; i < len(intfList); i++ {
		intf, err = net.InterfaceByName(intfList[i]) // bond0, eth0
		if err != nil {
			continue
		}
		break
	}
	// no intf found
	if intf == nil && err != nil {
		return "", err
	}

	// no ip found
	addrList := AddrList(intf)
	if len(addrList) == 0 {
		return "", errors.New("no ip found")
	}

	addrList = sortAddrList(addrList)
	return addrList[0].String(), nil
}

func AddrList(intf *net.Interface) []netip.Addr {

	addrs, _ := intf.Addrs()
	addrList := make([]netip.Addr, 0)
	// transfer net.Addr to netip.Addr
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			addr, ok := netip.AddrFromSlice(ipNet.IP.To4())
			if !ok {
				fmt.Println("AddrFromSlice failed.")
				continue
			}
			addrList = append(addrList, addr)
		}
	}
	return addrList
}

func sortAddrList(addrList []netip.Addr) []netip.Addr {
	sort.Slice(addrList, func(i, j int) bool {
		return addrList[i].Less(addrList[j])
	})
	return addrList
}

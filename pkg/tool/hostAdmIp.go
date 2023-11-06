package tool

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sort"
)

func HostAdmIp(intfList []string) (string, error) {
	if intfList == nil {
		intfList = []string{"bond0", "eth0", "eth1"}
	}
	return hostAdmIp(intfList)
}

// 管理地址所在interface规则，en0
// 选取最小的ip
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
	// no ip found
	if len(addrList) == 0 {
		return "", errors.New("no ip found")
	}
	addrList = sortAddrList(addrList)
	return addrList[0].String(), nil
}

func sortAddrList(addrList []netip.Addr) []netip.Addr {
	sort.Slice(addrList, func(i, j int) bool {
		return addrList[i].Less(addrList[j])
	})
	return addrList
}

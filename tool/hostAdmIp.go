// Copyright 2023 Sun Quan
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

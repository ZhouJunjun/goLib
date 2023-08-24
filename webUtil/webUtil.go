package webUtil

import (
	"fmt"
	"net"
)

func GetLocalServerIp() (string, error) {

	if ifaceSlice, err := net.Interfaces(); err == nil {
		if ifaceSlice != nil {
			for _, iface := range ifaceSlice {
				if iface.Flags&net.FlagLoopback != 0 {
					continue // loopback interface
				}
				if iface.Flags&net.FlagUp == 0 {
					continue // interface down
				}

				if tmpAddrSlice, err := iface.Addrs(); err == nil && tmpAddrSlice != nil {
					for _, address := range tmpAddrSlice {
						if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
							if ipNet.IP.To4() != nil {
								return ipNet.IP.String(), nil
							}
						}
					}
				}
			}
		}
		return "", fmt.Errorf("ip not found")
	} else {
		return "", err
	}
}

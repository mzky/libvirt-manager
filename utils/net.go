package utils

import (
	"net"
	"strings"
)

func appendIPNet(slice []net.IPNet, element net.IPNet) []net.IPNet {
	if element.IP.IsLinkLocalUnicast() { // ignore link local IPv6 address like "fe80::x"
		return slice
	}

	return append(slice, element)
}

func GetLocalIpNets() (map[string][]net.IPNet, error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	returnMap := make(map[string][]net.IPNet)
	for _, iFace := range iFaces {
		if iFace.Flags&net.FlagUp == 0 { // Ignore down adapter
			continue
		}

		if iFace.Flags&net.FlagLoopback == net.FlagLoopback { // Ignore loop back adapter
			continue
		}

		// if adapter name start with lo vir or tun(defined by kIgnoreAdapterPrefixes), we ignore it
		nameHasIgnorePrefix := false
		for _, ignoreName := range []string{"lo", "tun", "vir"} {
			if strings.HasPrefix(iFace.Name, ignoreName) {
				nameHasIgnorePrefix = true
				break
			}
		}
		if nameHasIgnorePrefix {
			continue
		}

		addrs, err := iFace.Addrs()
		if err != nil {
			continue
		}

		ipNets := make([]net.IPNet, 0)
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPAddr:
				ipNets = appendIPNet(ipNets, net.IPNet{IP: v.IP, Mask: v.IP.DefaultMask()})
			case *net.IPNet:
				ipNets = appendIPNet(ipNets, *v)
			}
		}
		returnMap[iFace.Name] = ipNets
	}

	return returnMap, nil
}

func GetLocalIPList() ([]string, error) {
	ipArray := make([]string, 0)
	ipMap, err := GetLocalIpNets()
	if err != nil {
		return nil, err
	}

	for _, ipNets := range ipMap {
		for _, ipNet := range ipNets {
			ipArray = append(ipArray, ipNet.IP.String())
		}
	}

	return ipArray, nil
}

// Contains 查找数组中是否包含字符串
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

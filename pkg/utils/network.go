package utils

import (
	"net"
	"strings"
	"time"
)

func GetOutboundIP() string {
	// 先尝试UDP方法
	if ip := GetOutboundIPByUDP(); ip != "127.0.0.1" && !IsTestNetwork(ip) {
		return ip
	}
	// 回退到检查网络接口
	return GetOutboundIPByInterface()
}

func GetOutboundIPByUDP() string {
	conn, err := net.DialTimeout("udp", "8.8.8.8:80", 3*time.Second)
	if err != nil {
		return "127.0.0.1"
	}
	defer func() {
		_ = conn.Close()
	}()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func IsTestNetwork(ip string) bool {
	// 检查是否是测试网络地址
	testNetworks := []string{
		"198.18.0.0/15",   // RFC 2544 测试网络
		"192.0.2.0/24",    // RFC 3330 测试网络
		"198.51.100.0/24", // RFC 3330 测试网络
		"203.0.113.0/24",  // RFC 3330 测试网络
	}

	for _, network := range testNetworks {
		_, cidr, _ := net.ParseCIDR(network)
		if cidr != nil && cidr.Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}

func GetOutboundIPByInterface() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}

	for _, iface := range interfaces {
		// 跳过回环和未启用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 跳过虚拟接口（通常以utun、tun、tap开头）
		if strings.HasPrefix(iface.Name, "utun") ||
			strings.HasPrefix(iface.Name, "tun") ||
			strings.HasPrefix(iface.Name, "tap") {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					// 跳过测试网络
					if !IsTestNetwork(ip) {
						return ip
					}
				}
			}
		}
	}

	return "127.0.0.1"
}

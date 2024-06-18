package netx

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

var ips []string

func init() {
	ips = []string{
		"http://myexternalip.com/raw",
		"http://ip.cip.cc",
		"https://4.ipw.cn",
		"http://test.ipw.cn",
	}
}

// 获取外网IP
func GetIpOut() string {
	ipc := make(chan string)
	for i := 0; i < len(ips); i++ {
		go func(net_addr string) {
			ipc <- doWork(net_addr)
		}(ips[i])
	}
	select {
	case ip := <-ipc:
		return ip
	case <-time.After(10 * time.Second):
		return ""
	}
}

func doWork(net_addr string) string {
	client := http.Client{
		Timeout: time.Millisecond * 100,
	}
	resp, err := client.Get(net_addr)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	return string(content)
}

// 获取内网IP
func GetIpIn() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

func IpToInt(ip_str string) uint32 {
	return binary.BigEndian.Uint32(net.ParseIP(ip_str).To4())
}

func IntToIp(ip_int uint32) string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, ip_int)
	return net.IPv4(b[0], b[1], b[2], b[3]).String()
}

package osx

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

// 上传网速
func NetSpeed(speed chan float64) {
	// 首次获取网络接口信息
	lasteNetIo, err := net.IOCounters(true)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// 再次获取网络接口信息
		nowNetIo, err := net.IOCounters(true)
		if err != nil {
			panic(err)
		}
		// 计算并打印每个网卡的上传速率
		var uploadSpeedAll float64
		for _, lastStat := range lasteNetIo {
			finalStat, exists := findInterfaceByName(nowNetIo, lastStat.Name)
			if !exists {
				fmt.Printf("Network interface %s disappeared.\n", lastStat.Name)
				continue
			}
			uploadSpeed := calculateUploadSpeed(lastStat.BytesSent, finalStat.BytesSent, time.Second)
			fmt.Printf("Interface: %s, Upload Speed: %.2f KB/s\n", lastStat.Name, uploadSpeed/1024)
			uploadSpeedAll += uploadSpeed
		}
		lasteNetIo = nowNetIo
		speed <- uploadSpeedAll
	}
}

// 在网络接口列表中根据名称查找对应接口
func findInterfaceByName(interfaces []net.IOCountersStat, name string) (net.IOCountersStat, bool) {
	for _, iface := range interfaces {
		if iface.Name == name {
			return iface, true
		}
	}
	return net.IOCountersStat{}, false
}

// 计算上传速率
func calculateUploadSpeed(prevSent, currSent uint64, duration time.Duration) float64 {
	deltaBytes := currSent - prevSent
	return float64(deltaBytes) / float64(duration.Seconds())
}

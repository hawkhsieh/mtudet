package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	protocolICMP = 1
	timeoutSec   = 3
)

func main() {
	// 解析命令行參數
	host := flag.String("host", "", "目標主機名或IP地址 (必填)")
	minMTU := flag.Int("min", 68, "最小MTU值 (最小值為68)")
	maxMTU := flag.Int("max", 1500, "最大MTU值")
	flag.Parse()

	if *host == "" {
		fmt.Println("Please specify a target host, e.g.: -host example.com")
		flag.Usage()
		os.Exit(1)
	}

	if *minMTU < 68 {
		fmt.Println("Minimum MTU value cannot be less than 68")
		os.Exit(1)
	}

	if *maxMTU <= *minMTU {
		fmt.Println("Maximum MTU value must be greater than minimum MTU value")
		os.Exit(1)
	}

	// 解析目標主機
	ipAddr, err := net.ResolveIPAddr("ip4", *host)
	if err != nil {
		fmt.Printf("Unable to resolve host %s: %v\n", *host, err)
		os.Exit(1)
	}

	fmt.Printf("Starting MTU detection for %s (%s)...\n", *host, ipAddr.String())

	// 使用二分搜尋法找出MTU
	mtu := findMTU(ipAddr.String(), *minMTU, *maxMTU)
	fmt.Printf("\nDetected MTU value: %d\n", mtu)
}

func findMTU(target string, min, max int) int {
	// 二分搜尋法
	for min <= max {
		mid := min + (max-min)/2
		fmt.Printf("Testing MTU value: %d...", mid)

		if pingWithSize(target, mid) {
			fmt.Println("Success")
			min = mid + 1
		} else {
			fmt.Println("Failed")
			max = mid - 1
		}
	}

	return max
}

func pingWithSize(target string, size int) bool {
	// 建立ICMP連接
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Printf("Unable to create ICMP connection: %v\n", err)
		return false
	}
	defer conn.Close()

	// 設置超時
	conn.SetDeadline(time.Now().Add(timeoutSec * time.Second))

	// 建立ICMP回顯請求
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: createPayload(size - 28), // 減去IP頭部(20字節)和ICMP頭部(8字節)
		},
	}

	// 序列化ICMP訊息
	binMsg, err := msg.Marshal(nil)
	if err != nil {
		fmt.Printf("Unable to marshal ICMP message: %v\n", err)
		return false
	}

	// 發送ICMP封包
	_, err = conn.WriteTo(binMsg, &net.IPAddr{IP: net.ParseIP(target)})
	if err != nil {
		return false
	}

	// 接收回應
	reply := make([]byte, 1500)
	_, _, err = conn.ReadFrom(reply)
	if err != nil {
		return false
	}

	return true
}

// 創建特定大小的數據負載
func createPayload(size int) []byte {
	if size <= 0 {
		return []byte{}
	}
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	return payload
}

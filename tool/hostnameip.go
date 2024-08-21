package tool

import (
	"errors"
	"fmt"
	"net"
	"os"
)

func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		msg := fmt.Sprintf("get hostname failed! err: %v\n", err)
		return "", errors.New(msg)
	}
	return hostname, nil
}

func GetHostIp(etherNetInterface string) (string, error) {
	// 获取指定的网络接口对象
	ethIns, err := net.InterfaceByName(etherNetInterface)
	if err != nil {
		msg := fmt.Sprintf("get net.InterfaceByName failed! err: %v\n", err)
		return "", errors.New(msg)
	}

	// 获取网络接口的地址
	addrs, err := ethIns.Addrs()
	if err != nil {
		msg := fmt.Sprintf("get ethIns.Addrs failed! err: %v\n", err)
		return "", errors.New(msg)
	}

	// 遍历所有的地址，找到第一个非回环的IPv4地址
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip := v.IP
			if !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil
			}
		case *net.IPAddr:
			ip := v.IP
			if !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	// 没有找到符合条件的IP地址，返回错误
	return "", errors.New("get host ip failed! no matching ip. ")
}

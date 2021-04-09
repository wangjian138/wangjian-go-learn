package utils

import (
	"fmt"
	"net"
	"strings"
)

// GetFreeIP 获取可用IP
func GetFreeIP(targetAddr string) (ip string, err error) {
	scheme := "http"
	addrs := strings.SplitN(targetAddr, "://", 2)

	if len(addrs) == 2 {
		scheme = addrs[0]
		targetAddr = addrs[1]
	}

	addrs = strings.SplitN(targetAddr, ":", 2)
	if len(addrs) != 2 {
		if scheme == "http" {
			targetAddr += ":80"
		} else if scheme == "https" {
			targetAddr += ":443"
		}
	}

	var conn net.Conn
	conn, err = net.Dial("tcp", targetAddr)
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			err = fmt.Errorf("获取可用IP错误：目标地址[%s]不存在服务", targetAddr)

		}

		err = fmt.Errorf("获取可用IP错误：%s", err)
		return
	}
	defer conn.Close()

	ip = conn.LocalAddr().(*net.TCPAddr).IP.String()
	return
}

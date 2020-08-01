package myUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ShowJsonFormat(v interface{}) (string, error) {
	tmp, ok := v.(string)
	var buff []byte
	var err error
	if !ok {
		buff, ok = v.([]byte)
		tmp = string(buff)
	}
	if ok && len(tmp) > 0 && (tmp[:1] == "{" || tmp[:1] == "[") {
		buff = []byte(tmp)
	} else {
		buff, err = json.Marshal(v)
		if err != nil {
			return "", err
		}
	}

	var out bytes.Buffer
	err = json.Indent(&out, buff, "", "  ")

	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func GetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jumpIp := ""
	var ipList []string
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && !strings.Contains(ipnet.String(), "/16") {
				tmps := strings.Split(ipnet.IP.To4().String(), ".")
				if tmps[3] == "1" { // NOTE 跳过 x.x.x.1 的ip地址
					jumpIp = ipnet.IP.To4().String()
					continue
				}
				ipList = append(ipList, ipnet.IP.To4().String())
			}
		}
	}
	if len(ipList) == 1 {
		return ipList[0]
	} else if len(ipList) > 0 { // 多个ip
		return ipList[rand.Intn(len(ipList))] // 随机输出
	}
	return jumpIp // 只有跳过的一个ip 的时候输出
}

func GetIps() (ipList []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && !strings.Contains(ipnet.String(), "/16") {
				ipList = append(ipList, ipnet.IP.To4().String())
			}
		}
	}
	return
}

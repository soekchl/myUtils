package myUtils

import (
	"fmt"
	"net"
	"os"
	"strings"
	"bytes"
	"encoding/json"
)

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
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && !strings.Contains(ipnet.String(), "/16") {
				return ipnet.IP.To4().String()
			}
		}
	}
	return ""
}

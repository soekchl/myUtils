package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const ()

var (
	config_data        sync.Map
	config_branch_name = "dev" // local test prod
	file_name          = ""
)

func GetConfigString(key string) string {
	if temp, ok := config_data.Load(key); ok {
		if val, ok := temp.(string); ok {
			return os.ExpandEnv(val)
		}
	}
	return ""
}

func GetConfigBool(key string) (bool, error) {
	n := GetConfigString(key)
	if n == "" {
		return false, fmt.Errorf("Not Found Key!")
	}

	return strconv.ParseBool(n)
}

func GetConfigInt(key string) (int, error) {
	n := GetConfigString(key)
	if n == "" {
		return 0, fmt.Errorf("Not Found!")
	}

	return strconv.Atoi(n)
}

func GetConfigDefaultString(key, defaultStr string) string {
	if temp, ok := config_data.Load(key); ok {
		if val, ok := temp.(string); ok {
			return os.ExpandEnv(val)
		}
	}
	return defaultStr
}

func GetConfigDefaultBool(key string, defaultBool bool) bool {
	n := GetConfigString(key)
	if n == "" {
		return defaultBool
	}
	b, err := strconv.ParseBool(n)
	if err != nil {
		return defaultBool
	}
	return b
}

func GetConfigDefaultInt(key string, defaultInt int) int {
	n := GetConfigString(key)
	if n == "" {
		return defaultInt
	}
	i, err := strconv.Atoi(n)
	if err != nil {
		return defaultInt
	}
	return i
}

func ReLoadConfig() {
	InitConfig(config_branch_name, file_name)
}

// 配置文件名 和 分支名称
func InitConfig(path_file_name, branch_name string) {
	config_branch_name = branch_name
	file_name = path_file_name
	getNameConfigDatas()
}

func readFile() string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}()
	fi, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func getNameConfigDatas() {
	file_str := readFile()
	strs := strings.Split(file_str, "[")

	temp_str := ""
	l := len(config_branch_name)
	for _, v := range strs {
		if len(v) > l &&
			strings.Compare(config_branch_name, v[:len(config_branch_name)]) == 0 {
			temp_str = v[l+1:]
			break
		}
	}
	for _, v := range strings.Split(temp_str, "\n") {
		if len(v) < 1 || v[0] == '#' {
			continue
		}

		// 删除注释
		n := strings.IndexByte(v, '#')
		if n != -1 {
			v = v[:n-1]
		}

		n = strings.IndexByte(v, '=')
		if n == -1 {
			continue
		}
		k := strings.Replace(v[:n], "\t", "", -1)
		val := strings.Replace(v[n+1:], "\t", "", -1)
		k = strings.Replace(k, " ", "", -1)
		val = strings.Replace(val, " ", "", -1)
		config_data.Store(k, val)
	}
}

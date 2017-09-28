// myUtils project myUtils.go
package myUtils

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	LeaveDebug = iota
	LeaveInfo
	LeaveNotice
	LeaveWarning
	LeaveError
	LeaveNoShow

	max_buff_size = 65536
	max_file_size = 1024 * 1024 * 50 // 50M
)

var (
	file_log_name string
	dir_log_name  = "myLog"
	file_name     = ""
	file_log_flag = false
	show_leave    = LeaveDebug // 默认全输出

	log_buff         = bytes.NewBuffer(make([]byte, max_buff_size))
	out_put_log_time = time.Second / 2
	out_put_log_chan = make(chan string, 100)
	enter            = "\n"
	_file_format     string
)

// 设定显示log等级
func SetShowLeave(leave int) {
	switch leave {
	case LeaveInfo:
		show_leave = LeaveInfo
		break
	case LeaveNotice:
		show_leave = LeaveNotice
		break
	case LeaveWarning:
		show_leave = LeaveWarning
		break
	case LeaveError:
		show_leave = LeaveError
		break
	case LeaveNoShow:
		show_leave = LeaveNoShow
	default:
		show_leave = LeaveDebug
		break
	}
}

func init() {
	if runtime.GOOS == "windows" {
		enter = "\r\n"
	} else {
		enter = "\n"
	}

	_file_format = "%s\\%s_%s_%d.log"
	if runtime.GOOS != "windows" {
		_file_format = "%s/%s_%s_%d.log"
	}
}

func SetOutputFileLog(log_file_name string) {
	file_name = log_file_name
	dir_log_name = fmt.Sprintf("%s_log", file_name)
	checkFileSize()
	file_log_flag = true
	log_buff.Reset()
	go outPutLogLoop()
}

func checkFileSize() {
	// 判断是否存在  判断大小
	var file os.FileInfo
	var name string
	var err error

	for i := 0; ; i++ {
		name = fmt.Sprintf(_file_format, dir_log_name, time.Now().Format("20060102"), file_name, i)
		file, err = os.Stat(name)
		if err != nil {
			break
		}
		if file.Size() < int64(max_file_size) {
			break
		}
	}
	file_log_name = name
}

func SetOutPutLogIntervalTime(interval int64) {
	if interval < 1 {
		return
	}
	out_put_log_time = time.Duration(interval)
}

func Debug(v ...interface{}) {
	if show_leave <= LeaveDebug || file_log_flag {
		myLog("[D]", show_leave <= LeaveDebug, v...)
	}
}

func Info(v ...interface{}) {
	if show_leave <= LeaveInfo || file_log_flag {
		myLog("[I]", show_leave <= LeaveInfo, v...)
	}
}

func Notice(v ...interface{}) {
	if show_leave <= LeaveNotice || file_log_flag {
		myLog("[N]", show_leave <= LeaveNotice, v...)
	}
}

func Warn(v ...interface{}) {
	if show_leave <= LeaveWarning || file_log_flag {
		myLog("[W]", show_leave <= LeaveWarning, v...)
	}
}

func Error(v ...interface{}) {
	if show_leave <= LeaveError || file_log_flag {
		myLog("【E】", show_leave <= LeaveError, v...)
	}
}

func myLog(mark string, show bool, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	outstring := fmt.Sprintf("%s %s %-16s %v%s",
		time.Now().Format("2006/01/02 15:04:05"), mark, fmt.Sprintf("%s:%d", filename, line), fmt.Sprint(v...), enter)

	if show {
		fmt.Print(outstring)
	}
	if file_log_flag {
		out_put_log_chan <- outstring
	}
}

func outPutLogLoop() {
	for file_log_flag {
		select {
		case <-time.After(out_put_log_time):
			if log_buff.Len() > 0 { //	一定时间段内没有输出过log 并且有log就输出
				outputLog()
			}
		case buff, ok := <-out_put_log_chan:
			if ok {
				if log_buff.Len()+len(buff) > max_buff_size { // 当缓存 超过限定的时候 提前输出
					outputLog()
				}
				log_buff.Write([]byte(buff)) // 写入到缓冲区
			}
		}
	}
}

func outputLog() {
	if _, err := os.Stat(dir_log_name); err != nil {
		if err := os.Mkdir(dir_log_name, 0755); err != nil {
			fmt.Println(err, "Mkdir")
			return
		}
	}

	file, err := os.OpenFile(file_log_name, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		file, err = os.Create(file_log_name)
		if err != nil {
			fmt.Println("Error!!! file", err)
			return
		}
	}

	file.Write(log_buff.Bytes())
	log_buff.Reset()
	file.Close()
	checkFileSize()
}

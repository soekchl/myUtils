// myUtils project myUtils.go
package myUtils

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	LeaveDebug = iota
	LeaveInfo
	LeaveNotice
	LeaveWarning
	LeaveError
	LeaveNoShow
)

var (
	file_log_name string
	dir_log_name  = "myLog"
	file_log_flag = false
	show_leave    = LeaveDebug // 默认全输出

	log_buff         = bytes.NewBuffer(make([]byte, 65536))
	log_buff_mutex   sync.Mutex
	out_put_log_time = time.Second / 2
	out_log_chan     = make(chan bool, 2)
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

func SetOutputFileLog(file_name string) {
	dir_log_name = fmt.Sprintf("%s_log", file_name)
	file_log_name = fmt.Sprintf("%s\\%s_%s.log", dir_log_name, time.Now().Format("20060102"), file_name)
	file_log_flag = true
	go outPutLogLoop()
}

func SetOutPutLogIntervalTime(interval int64) {
	if interval < 1 {
		return
	}
	out_put_log_time = time.Duration(interval)
}

func NowOutLog() {
	defer func() {
		if err := recover(); err != nil {
			Error(err)
		}
	}()
	out_log_chan <- true
	<-out_log_chan // all buffer logs is out file done.
	print("ok")
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
		color.Set(color.FgYellow, color.Bold)
		myLog("[W]", show_leave <= LeaveWarning, v...)
		color.Unset()
	}
}

func Error(v ...interface{}) {
	if show_leave <= LeaveError || file_log_flag {
		color.Set(color.FgRed, color.Bold)
		myLog("【E】", show_leave <= LeaveError, v...)
		color.Unset()
	}
}

func myLog(mark string, show bool, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	outstring := fmt.Sprintf("%s %s %-16s %v\n",
		time.Now().Format("2006/01/02 15:04:05"), mark, fmt.Sprintf("%s:%d", filename, line), fmt.Sprint(v...))

	if show {
		fmt.Print(outstring)
	}
	if file_log_flag {
		addOutPutLog(outstring)
	}
}

func addOutPutLog(out string) {
	if runtime.GOOS == "windows" {
		out = out + "\r\n"
	} else {
		out = out + "\n"
	}

	log_buff_mutex.Lock()
	defer log_buff_mutex.Unlock()

	log_buff.WriteString(out)
}

func outPutLogLoop() {
	var ok bool
	for file_log_flag {
		select {
		case _, ok = <-out_log_chan:
			if ok && log_buff.Len() > 0 {
				outputLog()
			}
			out_log_chan <- true
			return
		case _, ok = <-time.After(out_put_log_time):
			if ok && log_buff.Len() > 0 {
				outputLog()
			}
		}
	}
}

func outputLog() {
	if _, err := os.Stat(dir_log_name); err != nil {
		if err := os.Mkdir(dir_log_name, 0644); err != nil {
			fmt.Println(err, "Mkdir")
			return
		}
	}

	file, err := os.OpenFile(file_log_name, os.O_APPEND, 0644)
	if err != nil {
		file, err = os.Create(file_log_name)
		if err != nil {
			fmt.Println("Error!!! file", err)
			return
		}
	}
	defer file.Close()

	log_buff_mutex.Lock()
	defer log_buff_mutex.Unlock()

	file.Write(log_buff.Bytes())
	log_buff.Reset()
}

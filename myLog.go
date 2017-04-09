// myUtils project myUtils.go
package myUtils

import (
	"fmt"
	"os"
	"path"
	"runtime"
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

	outDebug
	outInfo
	outNotice
	outWarning
	outError
)

var (
	file_log_name string
	dir_log_name  = "myLog"
	file_log_flag = false
	show_leave    = LeaveDebug // 默认全输出
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

func Debug(v ...interface{}) {
	if show_leave <= LeaveDebug || file_log_flag {
		color.Set(color.FgMagenta, color.Bold)
		defer color.Unset()
		myLog(outDebug, "[D]", show_leave <= LeaveDebug, v...)
	}
}

func Info(v ...interface{}) {
	if show_leave <= LeaveInfo || file_log_flag {
		color.Set(color.FgBlue, color.Bold)
		defer color.Unset()
		myLog(outInfo, "[I]", show_leave <= LeaveInfo, v...)
	}
}

func Notice(v ...interface{}) {
	if show_leave <= LeaveNotice || file_log_flag {
		color.Set(color.FgGreen, color.Bold)
		defer color.Unset()
		myLog(outNotice, "[N]", show_leave <= LeaveNotice, v...)
	}
}

func Warn(v ...interface{}) {
	if show_leave <= LeaveWarning || file_log_flag {
		color.Set(color.FgYellow, color.Bold)
		defer color.Unset()
		myLog(outWarning, "[W]", show_leave <= LeaveWarning, v...)
	}
}

func Error(v ...interface{}) {
	if show_leave <= LeaveError || file_log_flag {
		color.Set(color.FgRed, color.Bold)
		defer color.Unset()
		myLog(outError, "【E】", show_leave <= LeaveError, v...)
	}
}

func SetOutputFileLog(file_name string) {
	dir_log_name = fmt.Sprintf("%s_log", file_name)
	file_log_name = fmt.Sprintf("%s\\%s_%s.log", dir_log_name, time.Now().Format("20060102"), file_name)
	file_log_flag = true
}

func myLog(out_flag int, mark string, show bool, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	outstring := fmt.Sprintf("%s %s %s %d ---> %v\n",
		time.Now().Format("2006/01/02 15:04:05"), mark, filename, line, fmt.Sprint(v...))

	if show {
		fmt.Print(outstring)
	}
	if file_log_flag {
		go outputLog(outstring, out_flag)
	}
}

func outputLog(out string, out_flag int) {
	if runtime.GOOS == "windows" {
		out = out + "\r\n"
	} else {
		out = out + "\n"
	}

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

	file.Write([]byte(out))

	// 在输出到out_flag 这块
	file_name := ""
	switch out_flag {
	case outDebug:
		file_name = fmt.Sprintf("%s_D", file_log_name)
		break
	case outInfo:
		file_name = fmt.Sprintf("%s_I", file_log_name)
		break
	case outError:
		file_name = fmt.Sprintf("%s_E", file_log_name)
		break
	case outWarning:
		file_name = fmt.Sprintf("%s_W", file_log_name)
		break
	case outNotice:
		file_name = fmt.Sprintf("%s_N", file_log_name)
		break
	}
	file1, err := os.OpenFile(file_name, os.O_APPEND, 0644)
	if err != nil {
		file1, err = os.Create(file_name)
		if err != nil {
			fmt.Println("Error!!! file1", err)
			return
		}
	}
	defer file1.Close()

	file1.Write([]byte(out))
}

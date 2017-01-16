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
)

var (
	file_log_name = time.Now().Format("20060102") + ".log"
	file_log_flag = false
	show_leave    = LeaveDebug // 默认全输出
)

func SetShowLeave(leave int) {
	if leave < LeaveDebug || leave > LeaveError {
		return
	}
	show_leave = leave
}

func Debug(v ...interface{}) {
	show := false
	if show_leave <= LeaveDebug {
		show = true
	}
	if show || file_log_flag {
		color.Set(color.FgMagenta, color.Bold)
		defer color.Unset()
		myLog("[D]", show, v...)
	}
}

func Info(v ...interface{}) {
	show := false
	if show_leave <= LeaveInfo {
		show = true
	}
	if show || file_log_flag {
		color.Set(color.FgBlue, color.Bold)
		defer color.Unset()
		myLog("[I]", show, v...)
	}
}

func Notice(v ...interface{}) {
	show := false
	if show_leave <= LeaveNotice {
		show = true
	}
	if show || file_log_flag {
		color.Set(color.FgGreen, color.Bold)
		defer color.Unset()
		myLog("[N]", show, v...)
	}
}

func Warn(v ...interface{}) {
	show := false
	if show_leave <= LeaveWarning {
		show = true
	}
	if show || file_log_flag {
		color.Set(color.FgYellow, color.Bold)
		defer color.Unset()
		myLog("[W]", show, v...)
	}
}

func Error(v ...interface{}) {
	show := false
	if show_leave <= LeaveError {
		show = true
	}
	if show || file_log_flag {
		color.Set(color.FgRed, color.Bold)
		defer color.Unset()
		myLog("【E】", show, v...)
	}
}

func SetOutputFileLog(file_name string, file_output_flag bool) {
	file_log_name = fmt.Sprintf("%s_%s.log", time.Now().Format("20060102"), file_name)
	file_log_flag = file_output_flag
}

func myLog(mark string, show bool, v ...interface{}) {
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
		outputLog(outstring)
	}
}

func outputLog(out string) {
	out = out + "\r\n"
	file, err := os.OpenFile(file_log_name, os.O_APPEND, 0644)
	if err != nil {
		file, err = os.Create(file_log_name)
		if err != nil {
			fmt.Println("Error!!!", err)
			return
		}
	}
	defer file.Close()

	file.Write([]byte(out))
}

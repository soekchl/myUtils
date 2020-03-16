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
	_ = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarning
	LevelError
	LevelNoShow

	max_buff_size = 65536
	max_file_size = 1024 * 1024 * 50 // 50M
)

var (
	file_log_name string
	dir_log_name  = "myLog"
	file_name     = ""
	file_log_flag = false
	show_leave    = LevelDebug // 默认全输出
	out_put_leave = LevelDebug // 默认全输出

	log_buff         = bytes.NewBuffer(make([]byte, max_buff_size))
	out_put_log_time = time.Second / 2
	out_put_log_chan = make(chan string, 100)
	enter            = "\n"
	_file_format     string

	lastLivingMsgCount = 0
)

// 设定显示log等级
func SetShowLevel(leave int) {
	show_leave = getLevel(leave)
}

// 设定输出log等级
func SetOutPutLevel(leave int) {
	out_put_leave = getLevel(leave)
}

func getLevel(leave int) int {
	switch leave {
	case LevelInfo:
		return LevelInfo
	case LevelNotice:
		return LevelNotice
	case LevelWarning:
		return LevelWarning
	case LevelError:
		return LevelError
	case LevelNoShow:
		return LevelNoShow
	}
	return LevelDebug
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

func Debugf(format string, v ...interface{}) {
	if show_leave <= LevelDebug || (file_log_flag && out_put_leave <= LevelDebug) {
		myLog("[D]", show_leave <= LevelDebug, out_put_leave <= LevelDebug, fmt.Sprintf(format, v...))
	}
}

func Infof(format string, v ...interface{}) {
	if show_leave <= LevelInfo || (file_log_flag && out_put_leave <= LevelInfo) {
		myLog("[I]", show_leave <= LevelInfo, out_put_leave <= LevelInfo, fmt.Sprintf(format, v...))
	}
}

func Noticef(format string, v ...interface{}) {
	if show_leave <= LevelNotice || (file_log_flag && out_put_leave <= LevelNotice) {
		myLog("[N]", show_leave <= LevelNotice, out_put_leave <= LevelNotice, fmt.Sprintf(format, v...))
	}
}

func Warnf(format string, v ...interface{}) {
	if show_leave <= LevelWarning || (file_log_flag && out_put_leave <= LevelWarning) {
		myLog("[W]", show_leave <= LevelWarning, out_put_leave <= LevelWarning, fmt.Sprintf(format, v...))
	}
}

func Errorf(format string, v ...interface{}) {
	if show_leave <= LevelError || (file_log_flag && out_put_leave <= LevelError) {
		myLog("【E】", show_leave <= LevelError, out_put_leave <= LevelError, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	if show_leave <= LevelDebug || (file_log_flag && out_put_leave <= LevelDebug) {
		myLog("[D]", show_leave <= LevelDebug, out_put_leave <= LevelDebug, v...)
	}
}

func Info(v ...interface{}) {
	if show_leave <= LevelInfo || (file_log_flag && out_put_leave <= LevelInfo) {
		myLog("[I]", show_leave <= LevelInfo, out_put_leave <= LevelInfo, v...)
	}
}

func Notice(v ...interface{}) {
	if show_leave <= LevelNotice || (file_log_flag && out_put_leave <= LevelNotice) {
		myLog("[N]", show_leave <= LevelNotice, out_put_leave <= LevelNotice, v...)
	}
}

func Warn(v ...interface{}) {
	if show_leave <= LevelWarning || (file_log_flag && out_put_leave <= LevelWarning) {
		myLog("[W]", show_leave <= LevelWarning, out_put_leave <= LevelWarning, v...)
	}
}

func Error(v ...interface{}) {
	if show_leave <= LevelError || (file_log_flag && out_put_leave <= LevelError) {
		myLog("【E】", show_leave <= LevelError, out_put_leave <= LevelError, v...)
	}
}

func LiveMsg(v ...interface{}) {
	//	lastLivingMsgCount
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	outstring := fmt.Sprintf("%s %-16s %v",
		time.Now().Format("2006/01/02 15:04:05"), fmt.Sprintf("%s:%d", filename, line), fmt.Sprint(v...))

	addMsg := ""
	for i := 0; i < lastLivingMsgCount; i++ {
		addMsg = fmt.Sprint(addMsg, "\b")
	}
	fmt.Print(addMsg, outstring)

	lastLivingMsgCount = len(outstring)
}

func myLog(mark string, show bool, out_put bool, v ...interface{}) {
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
	if file_log_flag && out_put {
		out_put_log_chan <- outstring
	}
}

func outPutLogLoop() {
	t := time.Now().UnixNano() // 最后一次输出log时间
	for file_log_flag {
		select {
		case <-time.After(out_put_log_time):
			if log_buff.Len() > 0 { //	等待后续log到一定时间 以后输出log
				outputLog()
				t = time.Now().UnixNano()
			}
		case buff, ok := <-out_put_log_chan:
			if ok {
				if log_buff.Len()+len(buff) > max_buff_size { // 当缓存 超过限定的时候 提前输出
					outputLog()
					t = time.Now().UnixNano()
				}
				log_buff.Write([]byte(buff)) // 写入到缓冲区
			}
		}

		// 当log 一定时间段内没有输出就输出一次log
		if log_buff.Len() > 0 && (time.Now().UnixNano()-t) > int64(out_put_log_time) {
			outputLog()
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

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
	out_put_log_time = int64(time.Second / 3)
	out_put_log_chan = make(chan string, 128)
	out_put_now_log  = make(chan struct{}, 10)
	enter            = "\n"
	_file_format     string

	lastLivingMsgCount = 0
	logOutFmt          string   // = "%s %-7s %-16s %s%s"
	outMarks           []string //     = []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "NoShow"}
)

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
	SetOutMark(1) // NOTE 默认第一种
}

// 设定显示log等级
func SetShowLevel(leave int) {
	show_leave = getLevel(leave)
}

// 设定输出log等级
func SetOutPutLevel(leave int) {
	out_put_leave = getLevel(leave)
}

func SetOutMark(mode int) string {
	if mode == 1 {
		outMarks = []string{"DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "NoShow"}
		logOutFmt = "%s %-7s %-16s %s%s"
	} else {
		outMarks = []string{"[D]", "[I]", "[N]", "[W]", "【E】", "NoShow"}
		logOutFmt = "%s %s %-16s %s%s"
	}
	return GetAllLevel()
}

func GetAllLevel() string {
	str := ""
	for k, _ := range outMarks {
		str += fmt.Sprintf("%v-%v  ", k+1, outMarks[k])
	}
	return str
}

func GetShowLevel() (int, string) {
	switch show_leave {
	case LevelDebug:
		return LevelDebug, outMarks[LevelDebug-1]
	case LevelInfo:
		return LevelInfo, outMarks[LevelInfo-1]
	case LevelNotice:
		return LevelNotice, outMarks[LevelNotice-1]
	case LevelWarning:
		return LevelWarning, outMarks[LevelWarning-1]
	case LevelError:
		return LevelError, outMarks[LevelError-1]
	case LevelNoShow:
		return LevelNoShow, outMarks[LevelNoShow-1]
	}
	return 0, ""
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

func getLevelString(leave int) string {
	switch leave {
	case LevelInfo:
		return outMarks[LevelInfo-1]
	case LevelNotice:
		return outMarks[LevelNotice-1]
	case LevelWarning:
		return outMarks[LevelWarning-1]
	case LevelError:
		return outMarks[LevelError-1]
	case LevelNoShow:
		return ""
	}
	return outMarks[LevelDebug-1]
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
	if interval < int64(time.Second)/100 {
		out_put_log_time = int64(time.Second) / 100
		return
	}
	out_put_log_time = interval
}

func Debugf(format string, v ...interface{}) {
	if show_leave <= LevelDebug || (file_log_flag && out_put_leave <= LevelDebug) {
		myLog(LevelDebug, show_leave <= LevelDebug, out_put_leave <= LevelDebug, fmt.Sprintf(format, v...))
	}
}

func Infof(format string, v ...interface{}) {
	if show_leave <= LevelInfo || (file_log_flag && out_put_leave <= LevelInfo) {
		myLog(LevelInfo, show_leave <= LevelInfo, out_put_leave <= LevelInfo, fmt.Sprintf(format, v...))
	}
}

func Noticef(format string, v ...interface{}) {
	if show_leave <= LevelNotice || (file_log_flag && out_put_leave <= LevelNotice) {
		myLog(LevelNotice, show_leave <= LevelNotice, out_put_leave <= LevelNotice, fmt.Sprintf(format, v...))
	}
}

func Warnf(format string, v ...interface{}) {
	if show_leave <= LevelWarning || (file_log_flag && out_put_leave <= LevelWarning) {
		myLog(LevelWarning, show_leave <= LevelWarning, out_put_leave <= LevelWarning, fmt.Sprintf(format, v...))
	}
}

func Errorf(format string, v ...interface{}) {
	if show_leave <= LevelError || (file_log_flag && out_put_leave <= LevelError) {
		myLog(LevelError, show_leave <= LevelError, out_put_leave <= LevelError, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	if show_leave <= LevelDebug || (file_log_flag && out_put_leave <= LevelDebug) {
		myLog(LevelDebug, show_leave <= LevelDebug, out_put_leave <= LevelDebug, fmt.Sprint(v...))
	}
}

func Info(v ...interface{}) {
	if show_leave <= LevelInfo || (file_log_flag && out_put_leave <= LevelInfo) {
		myLog(LevelInfo, show_leave <= LevelInfo, out_put_leave <= LevelInfo, fmt.Sprint(v...))
	}
}

func Notice(v ...interface{}) {
	if show_leave <= LevelNotice || (file_log_flag && out_put_leave <= LevelNotice) {
		myLog(LevelNotice, show_leave <= LevelNotice, out_put_leave <= LevelNotice, fmt.Sprint(v...))
	}
}

func Warn(v ...interface{}) {
	if show_leave <= LevelWarning || (file_log_flag && out_put_leave <= LevelWarning) {
		myLog(LevelWarning, show_leave <= LevelWarning, out_put_leave <= LevelWarning, fmt.Sprint(v...))
	}
}

func Error(v ...interface{}) {
	if show_leave <= LevelError || (file_log_flag && out_put_leave <= LevelError) {
		myLog(LevelError, show_leave <= LevelError, out_put_leave <= LevelError, fmt.Sprint(v...))
	}
}

// 每次输出都把上一次的结果清除
func LiveMsg(v ...interface{}) {
	//	lastLivingMsgCount
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	outstring := fmt.Sprintf("%s %-16s %v",
		time.Now().Format("2006/01/02 15:04:05"),
		fmt.Sprintf("%s:%d", filename, line),
		fmt.Sprint(v...),
	)

	addMsg := ""
	for i := 0; i < lastLivingMsgCount; i++ {
		addMsg = fmt.Sprint(addMsg, "\b")
	}
	fmt.Print(addMsg, outstring)

	lastLivingMsgCount = len(outstring)
}

func myLog(level int, show bool, out_put bool, v string) {
	if !out_put && !show {
		return
	}

	mark := getLevelString(level)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)

	if show {
		outFd := os.Stdout
		if level == LevelError {
			outFd = os.Stderr
		}

		fmt.Fprintf(outFd, logOutFmt,
			time.Now().Format("2006/01/02 15:04:05"),
			mark,
			fmt.Sprintf("%s:%d", filename, line),
			v,
			enter,
		)
	}
	if file_log_flag && out_put {
		out_put_log_chan <- fmt.Sprintf(logOutFmt, // "%s %-7s %-16s %s%s"
			time.Now().Format("2006/01/02 15:04:05"),
			mark,
			fmt.Sprintf("%s:%d", filename, line),
			v,
			enter,
		)
	}
}

func outPutLogLoop() {
	lastOutPutTime := time.Now().UnixNano() // 最后一次输出log时间
	outPutLeaveTime := out_put_log_time     // log输出剩余时间
	getOutPutLeaveTime := func() {
		if time.Now().UnixNano()-lastOutPutTime > out_put_log_time {
			outPutLeaveTime = 0
		} else {
			outPutLeaveTime = time.Now().UnixNano() - lastOutPutTime
		}
		lastOutPutTime = time.Now().UnixNano()
	}

	for file_log_flag {
		select {
		case <-time.After(time.Duration(outPutLeaveTime)):
			if log_buff.Len() > 0 { //	等待后续log到一定时间 以后输出log
				outputLog()
				getOutPutLeaveTime()
			}
			outPutLeaveTime = out_put_log_time // 重置时间
		case buff, ok := <-out_put_log_chan:
			if ok {
				if log_buff.Len()+len(buff) > max_buff_size { // 当缓存 超过限定的时候 提前输出
					outputLog()
					getOutPutLeaveTime()
					outPutLeaveTime = out_put_log_time // 重置时间
				}
				log_buff.Write([]byte(buff)) // 写入到缓冲区
			}
		case <-out_put_now_log:
			if log_buff.Len() > 0 {
				outputLog()
				getOutPutLeaveTime()
				outPutLeaveTime = out_put_log_time // 重置时间
			}
		}
	}
}

func Flush() {
	out_put_now_log <- struct{}{}
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
	defer file.Close()

	file.Write(log_buff.Bytes())
	log_buff.Reset()
	checkFileSize()
}

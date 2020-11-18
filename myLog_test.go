package myUtils

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) {
	err := os.RemoveAll("./test_log")
	if err != nil {
		panic(err)
	}
	SetOutputFileLog("test")
	defer Flush()

	f := func(level int, strErr, strWarn, strNotice, strInfo, strDebug string) {
		SetShowLevel(level)
		Error(strErr, "输出 Error")
		Warn(strWarn, "输出 Warn")
		Notice(strNotice, "输出 Notice")
		Info(strInfo, "输出 Info")
		Debug(strDebug, "输出 Debug")
	}

	f(LevelNoShow, "禁止", "禁止", "禁止", "禁止", "禁止")
	f(LevelError, "", "禁止", "禁止", "禁止", "禁止")
	f(LevelWarning, "", "", "禁止", "禁止", "禁止")
	f(LevelNotice, "", "", "", "禁止", "禁止")
	f(LevelInfo, "", "", "", "", "禁止")
	f(LevelDebug, "", "", "", "", "")

	Errorf("%s", "testErrorf")
	Noticef("%s", "test Noticef")
	Warnf("%s", "test Warnf")
	Infof("%s", "test Infof")
	Debugf("%s%s%s", "test Debugf", "hello", "debug")
	t.Log("LevelDebug", LevelDebug)
	t.Log("SetOutMark(1)", SetOutMark(1))
	t.Log("SetOutMark(2)", SetOutMark(2))
	Errorf("%s", "testErrorf")
	Noticef("%s", "test Noticef")
	Warnf("%s", "test Warnf")
	Infof("%s", "test Infof")
	Debugf("%s%s%s", "test Debugf", "hello", "debug")
}

func TestOutPutLog(t *testing.T) {
	Notice("notice")
	Notice("notice")

	SetOutPutLogIntervalTime(int64(time.Second))
	err := os.RemoveAll("./testOutPut_log")
	if err != nil {
		panic(err)
	}
	SetOutputFileLog("testOutPut")
	defer Flush()
	SetShowLevel(LevelNoShow)

	buff := bytes.NewBuffer(make([]byte, max_buff_size))
	for i := 0; i < max_buff_size; i++ {
		buff.WriteByte(byte('a' + i))
	}
	for i := 0; i < 1000; i++ {
		Notice(buff.String())
	}

	dir, err := ioutil.ReadDir("testOutPut_log")
	if err != nil {
		panic(err)
	}
	if len(dir) <= 1 {
		panic("输出文档个数错误，理应大于一个 前一个50M左右")
	}

}

func BenchmarkError(b *testing.B) {
	SetShowLevel(LevelDebug)
	defer Flush()
	for i := 0; i < b.N; i++ {
		Error("test", i, i+1, i+2, i+3, i+4, i+5)
		Errorf("test %v %v %v %v %v %v", i, i+1, i+2, i+3, i+4, i+5)
		Notice("test", i, i+1, i+2, i+3, i+4, i+5)
		Noticef("test %v %v %v %v %v %v", i, i+1, i+2, i+3, i+4, i+5)
	}
}

func BenchmarkLiveMsg(b *testing.B) {
	// cmd run  go test -v
	for i := 0; i < b.N; i++ {
		LiveMsg("test", i)
		time.Sleep(time.Nanosecond * 10)
	}
}

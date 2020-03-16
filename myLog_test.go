package myUtils

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	SetOutputFileLog("test")
	//	Error("test")
	//	Notice("Notice")
	//	Info("info")
	//	Warn("warn")

	SetShowLevel(LevelNoShow)
	//	Warn("leave_warn")
	//	Error("leave_error")

	SetOutPutLevel(LevelDebug)

	st := time.Now()
	for i := 0; i < 100000; i++ {
		Error("test", i, i+1, i+2, i+3, i+4, i+5)
		//		time.Sleep(time.Millisecond)
	}
	SetShowLevel(LevelNotice)
	Notice(time.Since(st))

	SetShowLevel(LevelNoShow)
	Error("test")
	SetOutPutLevel(LevelWarning)
	Notice("Notice")
	Info("info")
	Warn("warn")

	Errorf("%s", "testErrorf")
	Noticef("%s", "test Noticef")
	Warnf("%s", "test Warnf")
	Infof("%s", "test Infof")
	Debugf("%s%s%s", "test Debugf", "hello", "debug")

	time.Sleep(time.Second)

	for n := 0; n < 1024; n++ {
		LiveMsg("test", n)
	}
}

func BenchmarkError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Error("test", i)
	}
}

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

	SetShowLeave(LeaveNoShow)
	//	Warn("leave_warn")
	//	Error("leave_error")

	SetOutPutLeave(LeaveDebug)

	st := time.Now()
	for i := 0; i < 100000; i++ {
		Error("test", i, i+1, i+2, i+3, i+4, i+5)
		time.Sleep(time.Millisecond * 50)
	}
	SetShowLeave(LeaveNotice)
	Notice(time.Since(st))

	SetShowLeave(LeaveNoShow)
	Error("test")
	SetOutPutLeave(LeaveWarning)
	Notice("Notice")
	Info("info")
	Warn("warn")

	time.Sleep(time.Second)
}

func BenchmarkError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Error("test", i)
	}
}

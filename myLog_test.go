package myUtils

import (
	"testing"
	"time"
)

func Test(t *testing.T) {

	defer NowOutLog() // 留存在缓存中的数据输出到文件中

	SetOutputFileLog("test")
	Error("test")
	Notice("Notice")
	Info("info")
	Warn("warn")

	SetShowLeave(LeaveError)
	Warn("leave_warn")
	Error("leave_error")

	st := time.Now()
	for i := 0; i < 10; i++ {
		Error("test", i)
	}
	SetShowLeave(LeaveNotice)
	Notice(time.Since(st))

	SetShowLeave(LeaveNoShow)
	Error("test")
	Notice("Notice")
	Info("info")
	Warn("warn")
}

func BenchmarkError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Error("test", i)
	}
}

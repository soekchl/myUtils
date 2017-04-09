package myUtils

import (
	"testing"
	"time"
)

func Test_Error(t *testing.T) {
	SetOutputFileLog("test")
	Error("test")
	Notice("Notice")
	Info("info")
	Warn("warn")

	SetShowLeave(LeaveError)
	Warn("leave_warn")
	Error("leave_error")

	st := time.Now()
	for i := 0; i < 1000; i++ { //8-9
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

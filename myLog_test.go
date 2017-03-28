package myUtils

import (
	"testing"
)

func Test_Error(t *testing.T) {
	SetOutputFileLog("test", true)
	Error("test")
	Notice("Notice")
	Info("info")
	Warn("warn")
}

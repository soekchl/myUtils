package moniter

import (
	"testing"
)

func Test(t *testing.T) {

	t.Log(GetHtml())

	t.Log(GetDiskInfo())

	t.Log(GetMemInfo())

	t.Log(GetUseCpuPercent())

}

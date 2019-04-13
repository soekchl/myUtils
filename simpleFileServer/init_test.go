package simpleFileSystem

import (
	"testing"
)

func Test(t *testing.T) {
	err := Start(":9090", ".", 10)
	if err != nil {
		t.Error(err)
	}
}

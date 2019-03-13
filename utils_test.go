package myUtils

import (
	"testing"
)

func TestShowJsonFormat(t *testing.T) {
	type Road struct {
		Name   string
		Number int
	}
	roads := []Road{
		{"Diamond Fork", 29},
		{"Sheep Creek", 51},
	}

	t.Log(ShowJsonFormat(roads))

	t.Log(ShowJsonFormat("asdfasdf"))

	t.Log(ShowJsonFormat(123444))
	tt := make(map[string]string)
	tt["123"] = "321"
	t.Log(ShowJsonFormat(tt))

	t.Log(ShowJsonFormat(`{"Name": "Diamond Fork","Number": 29}`))
	t.Log(ShowJsonFormat([]byte(`{"Name": "Diamond Fork","Number": 29}`)))

}

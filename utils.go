package myUtils

import (
	"bytes"
	"encoding/json"
)

func ShowJsonFormat(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")

	if err != nil {
		return "", err
	}

	return out.String(), nil
}

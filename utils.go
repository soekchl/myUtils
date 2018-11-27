package myUtils

import (
	"bytes"
	"encoding/json"
)

func ShowJsonFormat(v interface{}) (string, error) {
	tmp, ok := v.(string)
	var buff []byte
	var err error
	if !ok {
		buff, ok = v.([]byte)
		tmp = string(buff)
	}
	if ok && len(tmp) > 0 && (tmp[:1] == "{" || tmp[:1] == "[") {
		buff = []byte(tmp)
	} else {
		buff, err = json.Marshal(v)
		if err != nil {
			return "", err
		}
	}

	var out bytes.Buffer
	err = json.Indent(&out, buff, "", "  ")

	if err != nil {
		return "", err
	}

	return out.String(), nil
}

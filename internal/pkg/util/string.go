package util

import "encoding/json"

func JsonEncode(data interface{}) (string, error) {
	b, e := json.Marshal(data)
	if nil != e {
		return ``, e
	}
	return string(b), e
}

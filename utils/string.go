package utils

import "encoding/json"

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func Itob(b int) bool {
	if b == 1 {
		return true
	}
	return false
}

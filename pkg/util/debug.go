package util

import (
	"encoding/json"
	"fmt"
)

func Debug(data any) {
	bytes, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(bytes))
}

func OutPut(data any) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

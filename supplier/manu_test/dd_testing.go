package main

import (
	"encoding/json"
	"lt-test/supplier/tools"
)

func main() {

	subContent := make(map[string]string)
	msg := "testing ...."
	subContent["content"] = msg

	tmpMsg, _ := json.Marshal(subContent)

	toDd := make(map[string]string)
	toDd["msgtype"] = "text"
	toDd["text"] = string(tmpMsg)
	message, _ := json.Marshal(toDd)
	tools.DdTalk(message)

}

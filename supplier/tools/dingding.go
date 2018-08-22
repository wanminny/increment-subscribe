package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	DD_TOKEN = "eb823e18168ee4579ef5f3597835c3de676ccb1ef4b402bc730024256bff87d9"
)

//钉钉报警
func DdTalk(msg []byte) {

	//组装要发送的内容格式；钉钉需要
	subContent := make(map[string]string)
	content := string(msg)
	subContent["content"] = content

	tmpMsg, _ := json.Marshal(subContent)

	toDd := make(map[string]string)
	toDd["msgtype"] = "text"
	toDd["text"] = string(tmpMsg)
	message, _ := json.Marshal(toDd)

	body := bytes.NewBuffer(message)
	url := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", DD_TOKEN)
	fmt.Printf("%s", body)
	res, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Println(err)
	} else {
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println(err)
		}
		log.Println(fmt.Sprintf("钉钉接口返回消息：%s", result))
	}
}

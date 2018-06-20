package tools

import (
	"os"
	"log"
	"bufio"
	"strings"
	"io"
)


var (
	//日志文件
	LOG_TO_FILE_NAME = "supplier_log.txt"

)
func LogTofile(dumpLog string)  {

	file,err := os.OpenFile(LOG_TO_FILE_NAME,os.O_CREATE | os.O_APPEND | os.O_WRONLY,0755)
	if err != nil{
		log.Println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(strings.NewReader(dumpLog))
	io.Copy(file,reader)
}

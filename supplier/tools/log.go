package tools

import (
	"os"
	"log"
	"bufio"
	"strings"
	"io"
	"io/ioutil"
	"fmt"
)


var (
	//日志文件
	LOG_TO_FILE_NAME = "logs/supplier_log.txt"

)

type Position struct {
	FileName string
	Pos string
}

func LogTofile(dumpLog string)  {

	file,err := os.OpenFile(LOG_TO_FILE_NAME,os.O_CREATE | os.O_APPEND | os.O_WRONLY,0755)
	if err != nil{
		log.Println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(strings.NewReader(dumpLog))
	io.Copy(file,reader)
}


func SaveToFile(info string,fileName string)  {

	file,err := os.OpenFile(fileName,os.O_CREATE | os.O_APPEND | os.O_WRONLY,0755)
	if err != nil{
		log.Println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(strings.NewReader(info))
	io.Copy(file,reader)
}


func ReadFileLast(fileName string) (pos Position ,err error) {

	file,err := os.Open(fileName)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	stream,err := ioutil.ReadAll(file)
	if err != nil{
		log.Println(err)
	}
	if len(stream) > 1 {
		//fmt.Printf("%s",string(stream))
		info := strings.Split(string(stream),"\n")
		length := len(info)
		//log.Printf("%d",length)
		//log.Printf("%v\n %d\n %s\n %s\n%s\n",info,length,info[0],info[1])
		//xxx,xx,xxx
		for k,v := range info{
			if  k ==  length -2 {
				lineInfo :=strings.Split(v,",")
				fmt.Println(lineInfo)
				pos.FileName = string(lineInfo[1])
				pos.Pos = string(lineInfo[2])
			}
		}
	}
	return
}
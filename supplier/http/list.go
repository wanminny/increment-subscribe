package http

import (
	"net/http"
	"log"
	"fmt"
)

//文件目录服务
func content(writer http.ResponseWriter, request *http.Request) {

	fmt.Fprintln(writer,"show log service! access path /log")
}


func StartHttpService()  {

	http.HandleFunc("/",content)
	dirHandle :=http.FileServer(http.Dir("logs"))
	//注意是 /log/ 以及上面的真实目录logs; ./logs doesn't work !
	http.Handle("/log/",http.StripPrefix("/log/",dirHandle))

	log.Println("start web server :5000")
	log.Fatal(http.ListenAndServe(":5000",nil))
}

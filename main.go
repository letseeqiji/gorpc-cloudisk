package main

import (
	"fileserver/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadHandlerMsg)
	http.HandleFunc("/file/get", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.FileDownloadHandler)

	err := http.ListenAndServe(":8090", nil)

	if err != nil {
		fmt.Printf("fail to start server , err: ", err.Error())
	}
}

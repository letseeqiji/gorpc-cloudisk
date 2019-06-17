package handler

import (
	"encoding/json"
	"fileserver/meta"
	"fileserver/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 引入模板
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "服务器错误")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接受上传文件
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Fail get file , err : ", err.Error())
			return
		}
		defer file.Close()

		// 存储文件信息到内存中
		fileMeta := meta.FileMeta{
			Filename:head.Filename,
			Location:"./tmp/"+head.Filename,
			UploadAt:time.Now().Format("2006-01-02 15:04:05"),
		}

		// 创建本地文件接受文件流
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Fail create file , err : ", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize,err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Fail copy file , err : ", err.Error())
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UploadFileMeta(fileMeta)

		fmt.Println(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusOK)
	}
}

// 上传完成后跳转到这个路由  使用http.redirect
func UploadHandlerMsg(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "伤上传成功")
}

// localhost:8090/file/get?filehash=491772f88411d7b941fedcd84994cdd68a54536e
func FileQueryHandler(w http.ResponseWriter, r * http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")

	fileMeta := meta.GetFileMeta(fileSha1)

	fileMetaString, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "获取文件失败")
		return
	}

	io.WriteString(w, string(fileMetaString))
}

// localhost:8090/file/download?filehash=b318e95277dd2115e50304a9aad1b05973854efe
func FileDownloadHandler(w http.ResponseWriter, r * http.Request) {
	r.ParseForm()

	fileSha1 := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileSha1)

	// 这里的过程是unkonwn的之前
	file, err := os.Open(fileMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "获取文件失败")
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "获取文件失败")
		return
	}

	w.Header().Set("Content-Type", "applaction/octect-stream")
	w.Header().Set("Content-disposition", "attachment;filename=\""+fileMeta.Filename+"\"")
	w.Write(data)
}

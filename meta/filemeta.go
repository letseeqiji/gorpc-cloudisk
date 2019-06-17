package meta

// 文件信息
type FileMeta struct {
	FileSha1 string
	Filename string
	FileSize int64
	Location string
	UploadAt string
}

// 在内存中存储文件信息
var fileMetas map[string]FileMeta

func init()  {
	// 自动初始化变量
	fileMetas = make(map[string]FileMeta)
}

// 吧上传的文件存储到内存中
func UploadFileMeta(fmeata FileMeta) {
	fileMetas[fmeata.FileSha1] = fmeata
}

// 根据sha1获取对应文件
func GetFileMeta(fileSha1 string) FileMeta  {
	return fileMetas[fileSha1]
}

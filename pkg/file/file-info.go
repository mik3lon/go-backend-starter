package file

type FileInfo struct {
	Filename    string
	ContentType string
	Size        int64
	Content     []byte
}

func NewFileInfo(filename string, contentType string, size int64, content []byte) *FileInfo {
	return &FileInfo{Filename: filename, ContentType: contentType, Size: size, Content: content}
}

package file

type UploadFiles []*UploadFile

type UploadFile struct {
	Url         string `json:"Url"`
	ContentType string `json:"ContentType"`
	Size        int64  `json:"Size"`
	Name        string `json:"Name"`
}

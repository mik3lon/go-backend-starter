package file

type ImageUploader interface {
	Upload(f FileInfo) (*UploadFile, error)
}

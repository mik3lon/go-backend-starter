package file

import "context"

type ImageUploader interface {
	Upload(ctx context.Context, f FileInfo) (*UploadFile, error)
}

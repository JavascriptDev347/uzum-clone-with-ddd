package media

import "errors"

var (
	ErrEmptyFile       = errors.New("file is empty")
	ErrUnsupportedType = errors.New("unsupported content type")
	ErrFileTooLarge    = errors.New("file exceeds size limit")
	ErrUploadFailed    = errors.New("failed to upload file")
	ErrDeleteFailed    = errors.New("failed to delete file")
)

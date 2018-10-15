package storage

type UploadType int

const (
	Image UploadType = iota
)

type StaticRepository interface {
	Upload(t UploadType, name string, data []byte) error
	Delete(t UploadType, name string) error
}

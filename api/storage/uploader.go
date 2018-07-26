package storage

type ImageUploader interface {
	Uploade(name string, data []byte) error
}

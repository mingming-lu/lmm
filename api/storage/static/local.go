package static

import (
	"bufio"
	"fmt"
	"lmm/api/storage"
	"os"
)

const (
	base  = "/static/"
	image = base + "/images/"
)

func toDir(t storage.UploadType) string {
	switch t {
	case storage.Image:
		return image
	default:
		panic(fmt.Sprintf("unknown upload type: %d", t))
	}
}

type localStaticRepository struct{}

func NewLocalStaticRepository() storage.StaticRepository {
	return &localStaticRepository{}
}

func (repo *localStaticRepository) Upload(t storage.UploadType, name string, data []byte) error {
	path := toDir(t)
	file, err := os.OpenFile(path+name, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	if _, err := w.Write(data); err != nil {
		return err
	}

	return w.Flush()
}

func (repo *localStaticRepository) Delete(t storage.UploadType, name string) error {
	return os.Remove(toDir(t) + name)
}

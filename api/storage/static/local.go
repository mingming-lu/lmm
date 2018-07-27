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

type LocalStaticRepository struct{}

func (repo *LocalStaticRepository) Upload(t storage.UploadType, name string, data []byte) error {
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

func (repo *LocalStaticRepository) Delete(t storage.UploadType, name string) error {
	return nil
}

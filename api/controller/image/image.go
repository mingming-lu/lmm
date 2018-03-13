package image

import (
	"io/ioutil"
	usecase "lmm/api/usecase/image"
	"lmm/api/usecase/user"
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"
)

func Upload(c *elesion.Context) {
	usr, err := user.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	c.Request.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	fhs := c.Request.MultipartForm.File["src"]

	type bytes = []byte
	bulkData := make([]bytes, 0)
	for _, fh := range fhs {
		contentType := fh.Header.Get("Content-Type")
		switch contentType {
		case "image/jpeg", "image/png":
		default:
			c.Status(http.StatusBadRequest).String("Invalid content type: " + contentType)
			return
		}

		f, err := fh.Open()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			c.Status(http.StatusBadRequest).String("Invalid input image").Error(err.Error())
			return
		}

		bulkData = append(bulkData, data)
	}

	err = usecase.Upload(usr.ID, bulkData)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusCreated).String("success")
}

func GetAllImages(c *elesion.Context) {
	userID, err := strconv.ParseInt(c.Params.ByName("user"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid user").Error(err.Error())
		return
	}

	images, err := usecase.AllImages(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(images)
}

func GetAllPhotos(c *elesion.Context) {
	userID, err := strconv.ParseInt(c.Params.ByName("user"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid user").Error(err.Error())
		return
	}

	photos, err := usecase.AllPhotos(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(photos)
}

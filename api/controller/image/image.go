package image

import (
	"encoding/json"
	"io/ioutil"
	account "lmm/api/context/account/appservice"
	model "lmm/api/domain/model/image"
	usecase "lmm/api/usecase/image"
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"
)

func Upload(c *elesion.Context) {
	usr, err := account.Verify(c.Request.Header.Get("Authorization"))
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

	if len(bulkData) == 0 {
		c.Status(http.StatusBadRequest).String("No file selected")
		return
	}

	err = usecase.Upload(usr.ID, bulkData)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusCreated).String("success")
}

func GetAllImages(c *elesion.Context) {
	userID, err := strconv.ParseUint(c.Params.ByName("user"), 10, 64)
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

func GetPhotos(c *elesion.Context) {
	queryParams := c.Query()
	photos, err := usecase.FetchPhotos(c.Params.ByName("user"), queryParams.Get("count"), queryParams.Get("page"))
	switch err {
	case nil:
		c.Status(http.StatusOK).JSON(photos)
	case usecase.ErrInvalidUserID, usecase.ErrInvalidCount, usecase.ErrInvalidPage:
		c.Status(http.StatusBadRequest).String(http.StatusText(http.StatusBadRequest)).Error(err.Error())
	default:
		c.Status(http.StatusInternalServerError).String(http.StatusText(http.StatusInternalServerError)).Error(err.Error())
	}
}

func PutPhoto(c *elesion.Context) {
	usr, err := account.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	m := model.Minimal{}
	if err := json.NewDecoder(c.Request.Body).Decode(&m); err != nil {
		c.Status(http.StatusBadRequest).String("Invalid body").Error(err.Error())
		return
	}

	if err := usecase.TurnOnPhotoSwitch(usr.ID, m.Name); err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

func RemovePhoto(c *elesion.Context) {
	usr, err := account.Verify(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String("Unauthorized, invalid token").Error(err.Error())
		return
	}

	m := model.Minimal{}
	if err := json.NewDecoder(c.Request.Body).Decode(&m); err != nil {
		c.Status(http.StatusBadRequest).String("Invalid body").Error(err.Error())
		return
	}

	if err := usecase.TurnOffPhotoSwitch(usr.ID, m.Name); err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).String("success")
}

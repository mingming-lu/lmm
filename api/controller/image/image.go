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

	file, handler, err := c.Request.FormFile("src")
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid image").Error(err.Error())
		return
	}
	defer file.Close()

	contentType := handler.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg", "image/png":
	default:
		c.Status(http.StatusBadRequest).String("Invalid content type: " + contentType)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid input image").Error(err.Error())
		return
	}

	err = usecase.Upload(usr.ID, c.Request.FormValue("type"), data)
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusCreated).String("success")
}

func GetPhotos(c *elesion.Context) {
	userID, err := strconv.ParseInt(c.Params.ByName("user"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("Invalid user").Error(err.Error())
		return
	}

	images, err := usecase.Find(userID, "photo")
	if err != nil {
		c.Status(http.StatusInternalServerError).String("Internal server error").Error(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(images)
}

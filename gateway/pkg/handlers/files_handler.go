package handlers

import (
	"gateway/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func (h *Handler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.JSONError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	src, err := file.Open()
	if err != nil {
		return utils.JSONError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	defer src.Close()

	db := h.MongoClient.Database("videos")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return utils.JSONError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	objectId, err := bucket.UploadFromStream(file.Filename, src)

	if err != nil {
		return utils.JSONError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, echo.Map{"objectId": objectId})
}

func (h *Handler) Download(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{})
}

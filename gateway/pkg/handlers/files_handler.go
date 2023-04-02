package handlers

import (
	"encoding/json"
	"gateway/pkg/services"
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

	access := c.Get("access").(services.UserMetadata)

	event := VideoUploadedEvent{
		ObjectId: objectId.String(),
		Email:    access.Email(),
	}
	data, _ := json.Marshal(event)

	h.QueueService.Publish(c.Request().Context(), "videos", data)

	return c.JSON(http.StatusOK, echo.Map{"objectId": objectId})
}

type VideoUploadedEvent struct {
	ObjectId string `json:"objectId"`
	Email    string `json:"email"`
}

func (h *Handler) Download(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{})
}

package handlers

import (
	"encoding/json"
	"gateway/internal/services"
	"gateway/internal/utils"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
)

type VideoUploadedEvent struct {
	ObjectId string `json:"objectId"`
	Email    string `json:"email"`
}

func (h *Handler) Upload(c echo.Context) error {
	fileName, src, err := getFileFromRequest(c)
	if err != nil {
		return utils.JSONError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	defer src.Close()

	objectId, err := h.StorageService.UploadFromStream("videos", fileName, src)
	if err != nil {
		return utils.JSONError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	access := c.Get("access").(services.UserMetadata)

	event := VideoUploadedEvent{
		ObjectId: objectId,
		Email:    access.Email(),
	}
	data, _ := json.Marshal(event)

	if err := h.QueueService.Publish(c.Request().Context(), "videos", data); err != nil {
		return utils.JSONError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, echo.Map{"objectId": objectId})
}

func getFileFromRequest(c echo.Context) (string, multipart.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", nil, err
	}
	src, err := file.Open()
	if err != nil {
		return "", nil, err
	}
	return file.Filename, src, nil
}

func (h *Handler) Download(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{})
}

package operations

import (
	"bytes"
	"context"
	"converter/config"
	"converter/internal/services"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VideoUploadedEvent struct {
	ObjectId string `json:"objectId"`
	Email    string `json:"email"`
}

type AudioExtractedEvent struct {
	ObjectId         string `json:"objectId"`
	OriginalFilename string `json:"originalFilename"`
	Email            string `json:"email"`
}

type TempFile struct {
	Name string
	File *os.File
}

type Processor struct {
	Storage services.StorageService
	Queue   services.QueueService
	Config  config.Config
}

func (p *Processor) ProcessMessage(msg services.Delivery) error {
	var event VideoUploadedEvent
	if err := json.Unmarshal(msg.Body(), &event); err != nil {
		return err
	}
	tmpFile, err := p.DownloadToTmpFile("videos", event.ObjectId)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.File.Name())

	convertedBuffer, err := runAudioExtraction(tmpFile.File.Name())
	if err != nil {
		return err
	}

	newFileName := strings.TrimSuffix(tmpFile.Name, filepath.Ext(tmpFile.Name)) + ".mp3"

	objectId, err := p.Storage.UploadFromStream("audios", newFileName, convertedBuffer)
	if err != nil {
		return err
	}

	audioExtractedEvent := AudioExtractedEvent{
		ObjectId:         objectId,
		OriginalFilename: tmpFile.Name,
		Email:            event.Email,
	}
	data, _ := json.Marshal(audioExtractedEvent)
	err = p.Queue.Publish(context.TODO(), p.Config.AudioQueue, data)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) DownloadToTmpFile(database, objectId string) (*TempFile, error) {
	file, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		return nil, err
	}
	dest, err := p.Storage.OpenDownloadStream(database, objectId)
	if err != nil {
		return nil, err
	}
	io.Copy(file, dest)

	return &TempFile{
		Name: dest.GetFile().Name,
		File: file,
	}, nil
}

func runAudioExtraction(filePath string) (io.Reader, error) {
	buffer := bytes.NewBuffer(nil)

	cmd := exec.Command("ffmpeg", "-i", filePath, "-vn", "-acodec", "libmp3lame", "-f", "mp3", "pipe:1")
	cmd.Stdout = buffer

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buffer, nil
}

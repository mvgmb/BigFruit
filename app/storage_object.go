package app

import (
	"os"

	storage_object "github.com/mvgmb/BigFruit/app/proto/storage_object"
)

type StorageObject struct {
	file *os.File
}

func NewStorageObject() *StorageObject {
	return &StorageObject{
		file: nil,
	}
}

// Upload writes a chunk of bytes into a file
func (e *StorageObject) Upload(uploadRequest *storage_object.UploadRequest) *storage_object.UploadResponse {
	filePath := uploadRequest.FilePath
	start := uploadRequest.Start
	bytes := uploadRequest.Bytes

	if e.file == nil || e.file.Name() != filePath {
		err := e.createFile(filePath)
		if err != nil {
			return &storage_object.UploadResponse{Error: err.Error()}
		}
	}
	_, err := e.file.WriteAt(bytes, start)
	if err != nil {
		return &storage_object.UploadResponse{Error: err.Error()}
	}
	return &storage_object.UploadResponse{}
}

// Download returns a chunk of bytes from a file
func (e *StorageObject) Download(downloadRequest *storage_object.DownloadRequest) *storage_object.DownloadResponse {
	filePath := downloadRequest.FilePath
	start := downloadRequest.Start
	offset := downloadRequest.Offset
	if e.file == nil || e.file.Name() != filePath {
		err := e.openFile(filePath)
		if err != nil {
			return &storage_object.DownloadResponse{Error: err.Error()}
		}
	}

	_, err := e.file.Seek(start, 0)
	if err != nil {
		return &storage_object.DownloadResponse{Error: err.Error()}
	}

	readBuffer := make([]byte, offset)
	noBytesRead, err := e.file.Read(readBuffer)
	if err != nil {
		return &storage_object.DownloadResponse{Error: err.Error()}
	}

	return &storage_object.DownloadResponse{
		Bytes: readBuffer[:noBytesRead],
	}
}

func (e *StorageObject) createFile(filePath string) error {
	if e.file != nil {
		e.file.Close()
		e.file = nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	e.file = file
	return nil
}

func (e *StorageObject) openFile(filePath string) error {
	if e.file != nil {
		e.file.Close()
		e.file = nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	e.file = file
	return nil
}

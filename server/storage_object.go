package server

import (
	"os"

	pb "github.com/mvgmb/BigFruit/proto"
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
func (e *StorageObject) Upload(uploadRequest *pb.StorageObjectUploadRequest) *pb.StorageObjectUploadResponse {
	filePath := uploadRequest.FilePath
	start := uploadRequest.Start
	bytes := uploadRequest.Bytes

	if e.file == nil || e.file.Name() != filePath {
		err := e.createFile(filePath)
		if err != nil {
			return &pb.StorageObjectUploadResponse{Error: err.Error()}
		}
	}
	_, err := e.file.WriteAt(bytes, start)
	if err != nil {
		return &pb.StorageObjectUploadResponse{Error: err.Error()}
	}
	return &pb.StorageObjectUploadResponse{}
}

// Download returns a chunk of bytes from a file
func (e *StorageObject) Download(downloadRequest *pb.StorageObjectDownloadRequest) *pb.StorageObjectDownloadResponse {
	filePath := downloadRequest.FilePath
	start := downloadRequest.Start
	offset := downloadRequest.Offset

	if e.file == nil || e.file.Name() != filePath {
		err := e.openFile(filePath)
		if err != nil {
			return &pb.StorageObjectDownloadResponse{Error: err.Error()}
		}
	}

	_, err := e.file.Seek(start, 0)
	if err != nil {
		return &pb.StorageObjectDownloadResponse{Error: err.Error()}
	}

	readBuffer := make([]byte, offset)
	noBytesRead, err := e.file.Read(readBuffer)
	if err != nil {
		return &pb.StorageObjectDownloadResponse{Error: err.Error()}
	}

	return &pb.StorageObjectDownloadResponse{
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

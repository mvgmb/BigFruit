package server

import (
	"os"
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
func (e *StorageObject) Upload(filePath string, start int64, bytes []byte) error {
	if e.file == nil || e.file.Name() != filePath {
		err := e.createFile(filePath)
		if err != nil {
			return err
		}
	}
	_, err := e.file.WriteAt(bytes, start)
	if err != nil {
		return err
	}
	return nil
}

// Download returns a chunk of bytes from a file
func (e *StorageObject) Download(filePath string, start int64, offset int) ([]byte, error) {
	if e.file == nil || e.file.Name() != filePath {
		err := e.openFile(filePath)
		if err != nil {
			return nil, err
		}
	}
	_, err := e.file.Seek(start, 0)
	if err != nil {
		return nil, err
	}

	readBuffer := make([]byte, offset)

	noBytesRead, err := e.file.Read(readBuffer)
	if err != nil {
		return []byte("EOF"), err
	}
	return readBuffer[:noBytesRead], nil
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

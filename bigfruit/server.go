package bigfruit

import (
	"fmt"
	"os"
)

const (
	noListeners = 20
	size        = 3
)

type Server struct {
	file       *os.File
	filesPaths map[string]string
}

func NewBigFruitServer() *Server {
	return &Server{
		file:       nil,
		filesPaths: make(map[string]string),
	}
}

func (e *Server) OpenFile(fileName string) error {
	if e.file != nil {
		e.file.Close()
		e.file = nil
	}

	file, err := os.Open(e.filesPaths[fileName])
	if err != nil {
		return err
	}

	e.file = file
	return nil
}

func (e *Server) CloseFile() error {
	if e.file == nil {
		return nil
	}

	err := e.file.Close()
	if err != nil {
		return err
	}

	e.file = nil
	return nil
}

func (e *Server) RegisterFilePath(fileName, filePath string) error {
	if _, ok := e.filesPaths[fileName]; ok {
		return fmt.Errorf("File %s already exists", fileName)
	}

	e.filesPaths[fileName] = filePath
	return nil
}

func (e *Server) UnregisterFilePath(fileName string) {
	delete(e.filesPaths, fileName)
}

// SeekBytes returns a chunk of bytes from a file
func (e *Server) SeekBytes(from int64) ([]byte, error) {
	if e.file == nil {
		return nil, fmt.Errorf("File is closed")
	}

	_, err := e.file.Seek(from, 0)
	if err != nil {
		return nil, err
	}

	readBuffer := make([]byte, size)

	noBytesRead, err := e.file.Read(readBuffer)
	if err != nil {
		return []byte("EOF"), err
	}

	return readBuffer[:noBytesRead], nil
}

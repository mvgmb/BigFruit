package bigfruit

import (
	"fmt"
	"os"
)

type Server struct {
	file      *os.File
	filesInfo map[string]Info
}

type Info struct {
	path string
	size int64
}

func NewBigFruitServer() *Server {
	return &Server{
		file:      nil,
		filesInfo: make(map[string]Info),
	}
}

func (e *Server) OpenFile(fileName string) error {
	if e.file != nil {
		e.file.Close()
		e.file = nil
	}

	file, err := os.Open(e.filesInfo[fileName].path)
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
	if _, ok := e.filesInfo[fileName]; ok {
		return fmt.Errorf("File %s already exists", fileName)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	// TODO delete this line
	fmt.Printf("%s is %d bytes long\n", fileName, fi.Size())

	e.filesInfo[fileName] = Info{
		path: filePath,
		size: fi.Size(),
	}

	return nil
}

func (e *Server) UnregisterFilePath(fileName string) {
	delete(e.filesInfo, fileName)
}

// SeekBytes returns a chunk of bytes from a file
func (e *Server) SeekBytes(from int64, size int) ([]byte, error) {
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

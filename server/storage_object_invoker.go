package server

import (
	"fmt"
	"log"
	"strconv"

	pb "github.com/mvgmb/BigFruit/proto"
)

type StorageObjectInvoker struct {
	storageObject *StorageObject
}

func NewStorageObjectInvoker() *StorageObjectInvoker {
	return &StorageObjectInvoker{
		storageObject: NewStorageObject(),
	}
}

func (e *StorageObjectInvoker) Invoke(req *pb.Message) ([]byte, error) {
	switch req.Key {
	case "StorageObject.Upload":
		if len(req.RawData) != 3 {
			return nil, fmt.Errorf(fmt.Sprintf("not enough arguments: needed 3, got: %d", len(req.RawData)))
		}
		log.Println(req)
		start, err := strconv.Atoi(string(req.RawData[1]))
		if err != nil {
			return nil, err
		}
		err = e.storageObject.Upload(string(req.RawData[0]), int64(start), req.RawData[2])
		if err != nil {
			return nil, err
		}
		return []byte("OK"), nil
	case "StorageObject.Download":
		if len(req.RawData) != 3 {
			return nil, fmt.Errorf(fmt.Sprintf("not enough arguments: needed 3, got: %d", len(req.RawData)))
		}
		start, err := strconv.ParseInt(string(req.RawData[1]), 10, 64)
		if err != nil {
			return nil, err
		}
		offset, err := strconv.Atoi(string(req.RawData[2]))
		if err != nil {
			return nil, err
		}
		bytes, err := e.storageObject.Download(string(req.RawData[0]), start, offset)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	default:
		return nil, fmt.Errorf("method not found")
	}
}

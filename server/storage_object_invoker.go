package server

import (
	"fmt"

	"github.com/golang/protobuf/proto"
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

func (e *StorageObjectInvoker) Invoke(messageType string, req proto.Message) (proto.Message, error) {
	switch messageType {
	case "*proto.StorageObjectUploadRequest":
		uploadRequest := req.(*pb.StorageObjectUploadRequest)
		return e.storageObject.Upload(uploadRequest), nil
	case "*proto.StorageObjectDownloadRequest":
		downloadRequest := req.(*pb.StorageObjectDownloadRequest)
		return e.storageObject.Download(downloadRequest), nil
	default:
		return nil, fmt.Errorf("procedure not found")
	}
}

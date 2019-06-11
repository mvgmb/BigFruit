package server

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	storage_object "github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/app/server/objects"
)

type StorageObjectInvoker struct {
	storageObject *objects.StorageObject
}

func NewStorageObjectInvoker() *StorageObjectInvoker {
	return &StorageObjectInvoker{
		storageObject: objects.NewStorageObject(),
	}
}

func (e *StorageObjectInvoker) Invoke(requestType string, req proto.Message) (proto.Message, error) {
	switch requestType {
	case "UploadRequest":
		uploadRequest := req.(*storage_object.UploadRequest)
		return e.storageObject.Upload(uploadRequest), nil
	case "DownloadRequest":
		downloadRequest := req.(*storage_object.DownloadRequest)
		return e.storageObject.Download(downloadRequest), nil
	default:
		return nil, fmt.Errorf("procedure not found")
	}
}

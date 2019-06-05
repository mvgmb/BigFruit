package server

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/app"
	storage_object "github.com/mvgmb/BigFruit/app/proto/storage_object"
)

type StorageObjectInvoker struct {
	storageObject *app.StorageObject
}

func NewStorageObjectInvoker() *StorageObjectInvoker {
	return &StorageObjectInvoker{
		storageObject: app.NewStorageObject(),
	}
}

func (e *StorageObjectInvoker) Invoke(messageType string, req proto.Message) (proto.Message, error) {
	switch messageType {
	case "*storage_object.UploadRequest":
		uploadRequest := req.(*storage_object.UploadRequest)
		return e.storageObject.Upload(uploadRequest), nil
	case "*storage_object.DownloadRequest":
		downloadRequest := req.(*storage_object.DownloadRequest)
		return e.storageObject.Download(downloadRequest), nil
	default:
		return nil, fmt.Errorf("procedure not found")
	}
}

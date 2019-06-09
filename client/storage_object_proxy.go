package client

import (
	"sync"

	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/util"
)

type StorageObjectProxy struct {
	requestor    *Requestor
	requestMutex *sync.Mutex
}

func NewStorageObjectProxy(options *util.Options) (*StorageObjectProxy, error) {
	requestor, err := NewRequestor(options)
	if err != nil {
		return nil, err
	}
	e := &StorageObjectProxy{
		requestor:    requestor,
		requestMutex: &sync.Mutex{},
	}
	return e, nil

}

// Upload writes a chunk of bytes into a file
func (e *StorageObjectProxy) Upload(req *storage_object.UploadRequest) (*storage_object.UploadResponse, error) {
	e.requestMutex.Lock()
	res, err := e.requestor.Invoke(req)
	if err != nil {
		return nil, err
	}
	e.requestMutex.Unlock()

	return res.(*storage_object.UploadResponse), nil
}

// Download returns a chunk of bytes from a file
func (e *StorageObjectProxy) Download(req *storage_object.DownloadRequest) (*storage_object.DownloadResponse, error) {
	e.requestMutex.Lock()
	res, err := e.requestor.Invoke(req)
	if err != nil {
		return nil, err
	}
	e.requestMutex.Unlock()

	return res.(*storage_object.DownloadResponse), nil
}

func (e *StorageObjectProxy) Close() {
	e.requestor.Close()
}

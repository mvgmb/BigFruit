package client

import (
	"sync"

	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/util"
)

type StorageObjectProxy struct {
	requestors []*Requestor
	roundRobin int
	robinMutex *sync.Mutex
}

func NewStorageObjectProxy(options []*util.Options) (*StorageObjectProxy, error) {
	var requestors []*Requestor
	for i := range options {
		requestor, err := NewRequestor(options[i])
		if err != nil {
			return nil, err
		}
		requestors = append(requestors, requestor)
	}
	e := &StorageObjectProxy{
		requestors: requestors,
		roundRobin: -1,
		robinMutex: &sync.Mutex{},
	}
	return e, nil
}

// Upload writes a chunk of bytes into a file
func (e *StorageObjectProxy) Upload(req *storage_object.UploadRequest) (*storage_object.UploadResponse, error) {
	roundRobin := e.getNextRobin()

	res, err := e.requestors[roundRobin].Invoke(req)
	if err != nil {
		return nil, err
	}
	return res.(*storage_object.UploadResponse), nil
}

// Download returns a chunk of bytes from a file
func (e *StorageObjectProxy) Download(req *storage_object.DownloadRequest) (*storage_object.DownloadResponse, error) {
	roundRobin := e.getNextRobin()

	res, err := e.requestors[roundRobin].Invoke(req)
	if err != nil {
		return nil, err
	}
	return res.(*storage_object.DownloadResponse), nil
}

func (e *StorageObjectProxy) getNextRobin() int {
	e.robinMutex.Lock()
	defer e.robinMutex.Unlock()

	e.roundRobin++
	if e.roundRobin >= len(e.requestors) {
		e.roundRobin = 0
	}
	return e.roundRobin
}

func (e *StorageObjectProxy) Close() {
	for i := range e.requestors {
		e.requestors[i].Close()
	}
}

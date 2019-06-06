package client

import (
	"sync"

	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/util"
)

type StorageObjectProxy struct {
	requestors   []*Requestor
	roundRobin   int
	robinMutex   *sync.Mutex
	robinMutexes []*sync.Mutex
}

var (
	lock   = 0
	unlock = 0
)

func NewStorageObjectProxy(options []*util.Options) (*StorageObjectProxy, error) {
	var requestors []*Requestor
	for i := range options {
		requestor, err := NewRequestor(options[i])
		if err != nil {
			return nil, err
		}
		requestors = append(requestors, requestor)
	}
	rm := make([]*sync.Mutex, len(options))
	for i := range rm {
		rm[i] = &sync.Mutex{}
	}
	e := &StorageObjectProxy{
		requestors:   requestors,
		roundRobin:   -1,
		robinMutex:   &sync.Mutex{},
		robinMutexes: rm,
	}
	return e, nil
}

// Upload writes a chunk of bytes into a file
func (e *StorageObjectProxy) Upload(req *storage_object.UploadRequest) (*storage_object.UploadResponse, error) {
	e.robinMutex.Lock()
	roundRobin := e.getNextRobin()
	e.robinMutex.Unlock()

	e.robinMutexes[roundRobin].Lock()
	res, err := e.requestors[roundRobin].Invoke(req)
	if err != nil {
		return nil, err
	}
	e.robinMutexes[roundRobin].Unlock()

	return res.(*storage_object.UploadResponse), nil
}

// Download returns a chunk of bytes from a file
func (e *StorageObjectProxy) Download(req *storage_object.DownloadRequest) (*storage_object.DownloadResponse, error) {
	e.robinMutex.Lock()
	roundRobin := e.getNextRobin()
	e.robinMutex.Unlock()

	e.robinMutexes[roundRobin].Lock()
	res, err := e.requestors[roundRobin].Invoke(req)
	if err != nil {
		return nil, err
	}
	e.robinMutexes[roundRobin].Unlock()

	return res.(*storage_object.DownloadResponse), nil
}

func (e *StorageObjectProxy) getNextRobin() int {
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

package client

import (
	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/util"
)

type StorageObjectProxy struct{}

func NewStorageObjectProxy() *StorageObjectProxy {
	return &StorageObjectProxy{}
}

// Upload writes a chunk of bytes into a file
func (e *StorageObjectProxy) Upload(reqApp chan *storage_object.UploadRequest, resApp chan *storage_object.UploadResponse) error {
	options := []*util.Options{
		&util.Options{
			Host:     "localhost",
			Port:     8080,
			Protocol: "tcp",
		},
	}

	reqProxy := make(chan proto.Message)
	resProxy := make(chan proto.Message)

	go func() {
		for {
			req, more := <-reqApp
			if more {
				reqProxy <- req
			} else {
				close(reqProxy)
				return
			}
		}

	}()

	done := make(chan bool)
	go func() {
		for {
			res, more := <-resProxy
			if more {
				resApp <- res.(*storage_object.UploadResponse)
			} else {
				close(resApp)
				done <- true
				return
			}
		}
	}()

	err := Call(reqProxy, resProxy, options, true)
	if err != nil {
		return err
	}

	<-done
	return nil
}

// Download returns a chunk of bytes from a file
func (e *StorageObjectProxy) Download(reqApp chan *storage_object.DownloadRequest, resApp chan *storage_object.DownloadResponse) error {
	options := []*util.Options{
		&util.Options{
			Host:     "localhost",
			Port:     8080,
			Protocol: "tcp",
		},
	}

	reqProxy := make(chan proto.Message)
	resProxy := make(chan proto.Message)

	go func() {
		for {
			req, more := <-reqApp
			if more {
				reqProxy <- req
			} else {
				close(reqProxy)
				return
			}
		}

	}()

	done := make(chan bool)
	go func() {
		for {
			res, more := <-resProxy
			if more {
				resApp <- res.(*storage_object.DownloadResponse)
			} else {
				close(resApp)
				done <- true
				return
			}
		}
	}()

	err := Call(reqProxy, resProxy, options, true)
	if err != nil {
		return err
	}

	<-done
	return nil
}

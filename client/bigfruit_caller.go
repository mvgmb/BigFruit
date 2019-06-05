package client

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/mvgmb/BigFruit/app/proto/storage_object"
	"github.com/mvgmb/BigFruit/util"
)

type BigFruit struct {
	proxy interface{}
}

const maxNoConcurrentRequestsPerServer = 3

func NewBigFruit() *BigFruit {
	return &BigFruit{}
}

func (e *BigFruit) Call(objectName, methodName string, options []*util.Options, replicate bool, reqCh, resCh chan proto.Message) error {
	// Initialize proxy
	switch objectName {
	case "StorageObject":
		proxy, err := NewStorageObjectProxy(options)
		if err != nil {
			return err
		}
		defer proxy.Close()
		e.proxy = proxy
	default:
		return fmt.Errorf("Object requested not found")
	}

	// Initialize internal channels
	var internal []chan proto.Message
	if replicate {
		internal = make([]chan proto.Message, len(options))
	} else {
		internal = make([]chan proto.Message, maxNoConcurrentRequestsPerServer*len(options))
	}

	errors := make([]chan error, len(internal))
	for i := 0; i < len(internal); i++ {
		internal[i] = make(chan proto.Message)
		errors[i] = make(chan error)
	}

	closed := false
	permission := make(chan bool)

	go func() {
		curID := 0
		requestorsRobin := 0
		noRoutinesPerReq := 1
		if replicate {
			noRoutinesPerReq = len(options)
		}

		for {
			request, more := <-reqCh
			if more {
				for i := 0; i < noRoutinesPerReq; i++ {
					<-permission

					go func(id, requestorIndex int) {
						index := id % len(internal)

						res, err := callObject(objectName, methodName, e.proxy, request)
						errors[index] <- err
						internal[index] <- res
					}(curID, requestorsRobin)

					curID++
					requestorsRobin++
					if requestorsRobin >= len(options) {
						requestorsRobin = 0
					}
				}
			} else {
				closed = true
				break
			}
		}
	}()

	// Initize as many go routines as possible
	for i := 0; i < len(internal); i++ {
		permission <- true
	}

	curID := 0
	// The idea is to consume and then initialize a new go routine
	for {
		i := curID % len(internal)

		err := <-errors[i]
		if err != nil {
			return err
		}

		res := <-internal[i]
		resCh <- res

		if closed {
			close(resCh)
			break
		}
		permission <- true
		curID++
	}

	return nil
}

func callObject(objectName, methodName string, proxy interface{}, req proto.Message) (proto.Message, error) {
	switch objectName {
	case "StorageObject":
		res, err := callStorageObjectMethod(methodName, proxy.(*StorageObjectProxy), req)
		return res, err
	default:
		return nil, fmt.Errorf("Object not found")
	}
}

func callStorageObjectMethod(methodName string, proxy *StorageObjectProxy, req proto.Message) (proto.Message, error) {
	switch methodName {
	case "Upload":
		res, err := proxy.Upload(req.(*storage_object.UploadRequest))
		return res, err
	case "Download":
		res, err := proxy.Download(req.(*storage_object.DownloadRequest))
		return res, err
	default:
		return nil, fmt.Errorf("Method not found")
	}
}

package client

import (
	"fmt"
	"strconv"

	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

type StorageObjectProxy struct {
	Requestor *Requestor
	Options   *util.Options
}

func NewStorageObjectProxy(requestor *Requestor, options *util.Options) *StorageObjectProxy {
	return &StorageObjectProxy{
		Requestor: requestor,
		Options:   options,
	}
}

// Upload writes a chunk of bytes into a file
func (e *StorageObjectProxy) Upload(filePath string, start int64, bytes *[]byte) error {
	req := util.NewMessage(200, "OK", "StorageObject.Upload",
		[]byte(filePath),
		[]byte(strconv.FormatInt(start, 10)),
		*bytes,
	)

	response, err := e.Requestor.Invoke(&req, e.Options)
	if err != nil {
		return err
	}
	res := response.(*pb.Message)

	if res.Status.Code != 200 {
		return fmt.Errorf(res.Status.Message, res.Key)
	}
	return nil
}

// Download returns a chunk of bytes from a file
func (e *StorageObjectProxy) Download(filePath string, start int64, offset int) ([]byte, error) {
	req := util.NewMessage(200, "OK", "StorageObject.Download",
		[]byte(filePath),
		[]byte(strconv.FormatInt(start, 10)),
		[]byte(strconv.Itoa(offset)),
	)

	response, err := e.Requestor.Invoke(&req, e.Options)
	if err != nil {
		return nil, err
	}
	res := response.(*pb.Message)

	if res.Status.Code != 200 {
		return nil, fmt.Errorf(res.Status.Message, res.Key)
	}
	return res.RawData[0], nil
}

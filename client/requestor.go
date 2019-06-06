package client

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

type Requestor struct {
	requestHandler *clientRequestHandler
}

func NewRequestor(options *util.Options) (*Requestor, error) {
	rh, err := newClientRequestHandler(options)
	if err != nil {
		return nil, err
	}
	e := &Requestor{requestHandler: rh}

	return e, nil
}

func (e *Requestor) Close() error {
	return e.requestHandler.close()
}

func (e *Requestor) Invoke(req proto.Message) (proto.Message, error) {
	bytes, err := util.WrapMessage(req)
	if err != nil {
		return nil, err
	}

	err = e.requestHandler.send(&bytes)
	if err != nil {
		return nil, err
	}

	bytes, err = e.requestHandler.receive()
	if err != nil {
		return nil, err
	}

	res, _, responseType, err := util.UnwrapMessage(bytes)
	if err != nil {
		return nil, err
	}

	if responseType == "Error" {
		return nil, fmt.Errorf((res.(*pb.Error)).Message)
	}

	return res, nil
}

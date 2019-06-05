package server

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

type Requestor struct {
	requestHandler *clientRequestHandler
	marshaller     *util.Marshaller
}

func NewRequestor() (*Requestor, error) {
	rh := newClientRequestHandler()

	marsh := util.NewMarshaller()

	e := &Requestor{
		requestHandler: rh,
		marshaller:     marsh,
	}

	return e, nil
}

func (e *Requestor) open(options *util.Options) error {
	return e.requestHandler.open(*options)
}

func (e *Requestor) close(options *util.Options) error {
	return e.requestHandler.close()
}

func (e *Requestor) invoke(req *proto.Message, options *util.Options) (proto.Message, error) {
	data, err := e.marshaller.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = e.requestHandler.send(&data)
	if err != nil {
		return nil, err
	}

	data, err = e.requestHandler.receive()
	if err != nil {
		return nil, err
	}

	res := pb.MessageWrapper{}

	err = e.marshaller.Unmarshal(&data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

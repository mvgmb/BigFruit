package client

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

type Requestor struct {
	requestHandler *clientRequestHandler
	marshaller     *util.Marshaller
}

func NewRequestor(options *util.Options) (*Requestor, error) {
	rh, err := newClientRequestHandler(options)
	if err != nil {
		return nil, err
	}
	marsh := util.NewMarshaller()

	e := &Requestor{
		requestHandler: rh,
		marshaller:     marsh,
	}

	return e, nil
}

func (e *Requestor) Close() error {
	return e.requestHandler.close()
}

func (e *Requestor) Invoke(req proto.Message) (proto.Message, error) {
	data, err := util.WrapMessage(req)
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
	wrapper := &pb.MessageWrapper{}

	err = e.marshaller.Unmarshal(&data, wrapper)
	if err != nil {
		return nil, err
	}

	res, err := util.UnwrapMessage(wrapper)
	if err != nil {
		return nil, err
	}

	return res, nil
}

package server

import (
	"github.com/mvgmb/BigFruit/proto/naming"
	"github.com/mvgmb/BigFruit/util"
)

var options = &util.Options{
	Host:     "localhost",
	Port:     1337,
	Protocol: "tcp",
}

func Bind(bindRequest *naming.BindRequest) (*naming.BindResponse, error) {
	requestor, err := NewRequestor(options)
	if err != nil {
		return nil, err
	}
	res, err := requestor.Invoke(bindRequest)
	if err != nil {
		return nil, err
	}
	requestor.Close()

	return res.(*naming.BindResponse), nil
}

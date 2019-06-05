package server

import (
	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

var options = &util.Options{
	Host:     "localhost",
	Port:     1337,
	Protocol: "tcp",
}

func Bind(bindRequest *pb.NamingServiceBindRequest) (*pb.NamingServiceBindResponse, error) {
	requestor, err := NewRequestor(options)
	if err != nil {
		return nil, err
	}
	res, err := requestor.Invoke(bindRequest)
	if err != nil {
		return nil, err
	}
	requestor.Close()

	return res.(*pb.NamingServiceBindResponse), nil
}

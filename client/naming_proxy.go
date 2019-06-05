package client

import (
	pb "github.com/mvgmb/BigFruit/proto"
	"github.com/mvgmb/BigFruit/util"
)

var lookupOptions = &util.Options{
	Host:     "localhost",
	Port:     1337,
	Protocol: "tcp",
}

func LookupMany(lookupManyRequest *pb.NamingServiceLookupManyRequest) (*pb.NamingServiceLookupManyResponse, error) {
	requestor, err := NewRequestor(lookupOptions)
	if err != nil {
		return nil, err
	}
	res, err := requestor.Invoke(lookupManyRequest)
	if err != nil {
		return nil, err
	}
	requestor.Close()

	return res.(*pb.NamingServiceLookupManyResponse), nil
}

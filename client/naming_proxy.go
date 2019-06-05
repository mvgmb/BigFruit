package client

import (
	"github.com/mvgmb/BigFruit/proto/naming"
	"github.com/mvgmb/BigFruit/util"
)

var lookupOptions = &util.Options{
	Host:     "localhost",
	Port:     1337,
	Protocol: "tcp",
}

func LookupMany(lookupManyRequest *naming.LookupManyRequest) (*naming.LookupManyResponse, error) {
	requestor, err := NewRequestor(lookupOptions)
	if err != nil {
		return nil, err
	}
	res, err := requestor.Invoke(lookupManyRequest)
	if err != nil {
		return nil, err
	}
	requestor.Close()

	return res.(*naming.LookupManyResponse), nil
}

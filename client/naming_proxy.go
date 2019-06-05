package client

import (
	"github.com/mvgmb/BigFruit/util"
)

var lookupOptions = &util.Options{
	Host:     "localhost",
	Port:     1337,
	Protocol: "tcp",
}

func (e *Requestor) Lookup(serviceName string) (*util.Options, error) {
	// req := util.NewMessage(200, "OK", "Lookup", []byte(serviceName))

	// err := e.Open(lookupOptions)
	// if err != nil {
	// 	return nil, err
	// }

	// result, err := e.Invoke(&req, lookupOptions)
	// if err != nil {
	// 	return nil, err
	// }

	// err = e.Close()
	// if err != nil {
	// 	return nil, err
	// }

	// res, ok := result.(*pb.MessageWrapper)
	// if !ok {
	// 	return nil, fmt.Errorf("Not a Message")
	// }

	// aor, err := util.StringToAOR(string(res.RawData[0]))
	// if err != nil {
	// 	return nil, err
	// }

	// options := &util.Options{
	// 	Host:     aor.Host,
	// 	Port:     aor.Port,
	// 	Protocol: "tcp",
	// }

	// return options, nil
	return nil, nil
}

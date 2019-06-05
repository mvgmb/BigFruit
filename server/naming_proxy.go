package server

import (
	"github.com/mvgmb/BigFruit/util"
)

var options = util.Options{
	Host:     "localhost",
	Port:     1337,
	Protocol: "tcp",
}

func Bind(aor *util.AOR) error {
	// requestor, err := NewRequestor()
	// if err != nil {
	// 	return err
	// }
	// req := util.NewMessageWrapper(200, "OK", "Bind", []byte(aor.String()))

	// err = requestor.open(&options)
	// if err != nil {
	// 	return err
	// }

	// _, err = requestor.invoke(&req, &options)
	// if err != nil {
	// 	return err
	// }

	// err = requestor.close(&options)
	// if err != nil {
	// 	return err
	// }

	return nil
}

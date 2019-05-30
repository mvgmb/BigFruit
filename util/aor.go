package util

import (
	"fmt"
	"strconv"
	"strings"
)

// AOR is the Absolute Object Reference
type AOR struct {
	Host string
	Port uint16
	ID   string
}

func (e *AOR) String() string {
	return fmt.Sprintf("%s:%d:%s", e.Host, e.Port, e.ID)
}

func StringToAOR(aor string) (*AOR, error) {
	split := strings.Split(aor, ":")

	num, err := strconv.ParseUint(split[1], 10, 16)
	if err != nil {
		return nil, err
	}

	e := AOR{
		Host: split[0],
		Port: uint16(num),
		ID:   split[2],
	}

	return &e, nil
}

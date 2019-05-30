package util

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/mvgmb/BigFruit/proto"
)

var (
	ErrUnknown            = NewMessage(000, "Unknown", "")
	ErrBadRequest         = NewMessage(400, "Bad Request", "")
	ErrUnauthorized       = NewMessage(401, "Unauthorized", "")
	ErrForbidden          = NewMessage(403, "Forbidden", "")
	ErrNotFound           = NewMessage(404, "Not found", "")
	ErrMethodNotAllowed   = NewMessage(405, "Method not allowed", "")
	ErrPayloadTooLarge    = NewMessage(413, "Payload too large!", "")
	ErrExpectationFailed  = NewMessage(417, "Expectation fail", "")
	ErrServiceUnavailable = NewMessage(503, "Service unavailable", "")
)

// Options defines the options values
type Options struct {
	Host     string
	Port     uint16
	Protocol string
}

// NewMessage creates a message
func NewMessage(statusCode uint64, statusMessage, key string, bytes ...[]byte) proto.Message {
	status := &pb.Status{
		Code:    statusCode,
		Message: statusMessage,
	}
	message := &pb.Message{
		Status:  status,
		Key:     key,
		RawData: bytes,
	}
	return message
}

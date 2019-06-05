package util

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	storage_object "github.com/mvgmb/BigFruit/app/proto/storage_object"
	pb "github.com/mvgmb/BigFruit/proto"
	naming "github.com/mvgmb/BigFruit/proto/naming"
)

var (
	// ErrUnknown            = NewMessage(000, "Unknown", "")
	ErrBadRequest = &pb.Error{Code: 400, Message: "Bad Request"}
	// ErrUnauthorized       = NewMessage(401, "Unauthorized", "")
	// ErrForbidden          = NewMessage(403, "Forbidden", "")
	// ErrNotFound           = NewMessage(404, "Not found", "")
	// ErrMethodNotAllowed   = NewMessage(405, "Method not allowed", "")
	// ErrPayloadTooLarge    = NewMessage(413, "Payload too large!", "")
	// ErrExpectationFailed  = NewMessage(417, "Expectation fail", "")
	// ErrServiceUnavailable = NewMessage(503, "Service unavailable", "")
)

type Options struct {
	Host     string
	Port     uint16
	Protocol string
}

var protoType = map[string]proto.Message{
	"*proto.Error": &pb.Error{},

	"*storage_object.UploadRequest":    &storage_object.UploadRequest{},
	"*storage_object.UploadResponse":   &storage_object.UploadResponse{},
	"*storage_object.DownloadRequest":  &storage_object.DownloadRequest{},
	"*storage_object.DownloadResponse": &storage_object.DownloadResponse{},

	"*naming.BindRequest":        &naming.BindRequest{},
	"*naming.BindResponse":       &naming.BindResponse{},
	"*naming.LookupRequest":      &naming.LookupRequest{},
	"*naming.LookupResponse":     &naming.LookupResponse{},
	"*naming.LookupManyRequest":  &naming.LookupManyRequest{},
	"*naming.LookupManyResponse": &naming.LookupManyResponse{},
	"*naming.LookupAllRequest":   &naming.LookupAllRequest{},
	"*naming.LookupAllResponse":  &naming.LookupAllResponse{},
}

func WrapMessage(message proto.Message) ([]byte, error) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	wrapper := &pb.MessageWrapper{
		Type:    reflect.TypeOf(message).String(),
		Message: bytes,
	}

	bytes, err = proto.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func UnwrapMessage(wrapper *pb.MessageWrapper) (proto.Message, error) {
	message := protoType[wrapper.Type]

	err := proto.Unmarshal(wrapper.Message, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

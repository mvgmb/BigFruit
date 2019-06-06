package util

import (
	"reflect"
	"strings"

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

func UnwrapMessage(bytes []byte) (proto.Message, string, string, error) {
	wrapper := &pb.MessageWrapper{}

	err := proto.Unmarshal(bytes, wrapper)
	if err != nil {
		return nil, "", "", err
	}

	message := protoType(wrapper.Type)

	err = proto.Unmarshal(wrapper.Message, message)
	if err != nil {
		return nil, "", "", err
	}

	protoType := strings.Split(wrapper.Type, ".")

	return message, protoType[0], protoType[1], nil
}

func protoType(messageType string) proto.Message {
	switch messageType {
	case "*proto.Error":
		return &pb.Error{}

	case "*storage_object.UploadRequest":
		return &storage_object.UploadRequest{}
	case "*storage_object.UploadResponse":
		return &storage_object.UploadResponse{}
	case "*storage_object.DownloadRequest":
		return &storage_object.DownloadRequest{}
	case "*storage_object.DownloadResponse":
		return &storage_object.DownloadResponse{}

	case "*naming.BindRequest":
		return &naming.BindRequest{}
	case "*naming.BindResponse":
		return &naming.BindResponse{}
	case "*naming.LookupRequest":
		return &naming.LookupRequest{}
	case "*naming.LookupResponse":
		return &naming.LookupResponse{}
	case "*naming.LookupManyRequest":
		return &naming.LookupManyRequest{}
	case "*naming.LookupManyResponse":
		return &naming.LookupManyResponse{}
	case "*naming.LookupAllRequest":
		return &naming.LookupAllRequest{}
	case "*naming.LookupAllResponse":
		return &naming.LookupAllResponse{}
	default:
		return nil
	}
}

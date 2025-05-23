// Based on https://github.com/trusch/grpc-proxy
// Copyright Michal Witkowski. Licensed under Apache2 license: https://github.com/trusch/grpc-proxy/blob/master/LICENSE.txt

package codec

import (
	"fmt"

	protoV1 "github.com/golang/protobuf/proto" //nolint:staticcheck
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

// Name is the name by which the proxy codec is registered in the encoding codec registry
// We have to say that we are the "proto" codec otherwise marshaling will fail.
const Name = "proto"

// Register manually registers the codec.
func Register() {
	encoding.RegisterCodec(codec())
}

// codec returns a proxying grpc.codec with the default protobuf codec as parent.
//
// See CodecWithParent.
func codec() encoding.Codec {
	// since we have registered the default codec by importing it,
	// we can fetch it from the registry and use it as our parent
	// and overwrite the existing codec in the registry
	return codecWithParent(&protoCodec{})
}

// CodecWithParent returns a proxying grpc.Codec with a user provided codec as parent.
//
// This codec is *crucial* to the functioning of the proxy. It allows the proxy server to be oblivious
// to the schema of the forwarded messages. It basically treats a gRPC message frame as raw bytes.
// However, if the server handler, or the client caller are not proxy-internal functions it will fall back
// to trying to decode the message using a fallback codec.
func codecWithParent(fallback encoding.Codec) encoding.Codec {
	return &Proxy{parentCodec: fallback}
}

// Proxy satisfies the encoding.Codec interface.
type Proxy struct {
	parentCodec encoding.Codec
}

// Frame holds the proxy transported data.
type Frame struct {
	payload []byte
}

// ProtoMessage tags a frame as valid proto message.
func (f *Frame) ProtoMessage() {
	// nop
}

// Marshal implements the encoding.Codec interface method.
func (p *Proxy) Marshal(v any) ([]byte, error) {
	out, ok := v.(*Frame)
	if !ok {
		return p.parentCodec.Marshal(v)
	}

	return out.payload, nil
}

// Unmarshal implements the encoding.Codec interface method.
func (p *Proxy) Unmarshal(data []byte, v any) error {
	dst, ok := v.(*Frame)
	if !ok {
		return p.parentCodec.Unmarshal(data, v)
	}
	dst.payload = data
	return nil
}

// Name implements the encoding.Codec interface method.
func (*Proxy) Name() string {
	return Name
}

// protoCodec is a Codec implementation with protobuf. It is the default rawCodec for gRPC.
type protoCodec struct{}

func (*protoCodec) Marshal(v any) ([]byte, error) {
	switch x := v.(type) {
	case proto.Message:
		return proto.Marshal(x)
	case protoV1.Message:
		return protoV1.Marshal(x)
	default:
		return nil, fmt.Errorf("failed to marshal: message is %T, want proto.Message", x)
	}
}

func (*protoCodec) Unmarshal(data []byte, v any) error {
	switch x := v.(type) {
	case proto.Message:
		return proto.Unmarshal(data, x)
	case protoV1.Message:
		return protoV1.Unmarshal(data, x)
	default:
		return fmt.Errorf("failed to unmarshal: message is %T, want proto.Message", x)
	}
}

func (*protoCodec) Name() string {
	return "proxy>proto"
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: prediction.proto

package prediction

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State string `protobuf:"bytes,1,opt,name=State,proto3" json:"State,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Limit int64  `protobuf:"varint,3,opt,name=Limit,proto3" json:"Limit,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prediction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_prediction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_prediction_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *Message) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Message) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type MessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code       int32   `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Prediction int64   `protobuf:"varint,2,opt,name=Prediction,proto3" json:"Prediction,omitempty"`
	Score      float32 `protobuf:"fixed32,3,opt,name=Score,proto3" json:"Score,omitempty"`
}

func (x *MessageResponse) Reset() {
	*x = MessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prediction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageResponse) ProtoMessage() {}

func (x *MessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_prediction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageResponse.ProtoReflect.Descriptor instead.
func (*MessageResponse) Descriptor() ([]byte, []int) {
	return file_prediction_proto_rawDescGZIP(), []int{1}
}

func (x *MessageResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *MessageResponse) GetPrediction() int64 {
	if x != nil {
		return x.Prediction
	}
	return 0
}

func (x *MessageResponse) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

var File_prediction_proto protoreflect.FileDescriptor

var file_prediction_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x49,
	0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x5b, 0x0a, 0x0f, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x32, 0xb1, 0x01, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x64, 0x69,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4d, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x53, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x13, 0x2e, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x76, 0x65, 0x72,
	0x73, 0x65, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x72, 0x50, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1b, 0x2e, 0x70, 0x72,
	0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x75, 0x73, 0x74, 0x68, 0x75, 0x6d,
	0x61, 0x6e, 0x7a, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_prediction_proto_rawDescOnce sync.Once
	file_prediction_proto_rawDescData = file_prediction_proto_rawDesc
)

func file_prediction_proto_rawDescGZIP() []byte {
	file_prediction_proto_rawDescOnce.Do(func() {
		file_prediction_proto_rawDescData = protoimpl.X.CompressGZIP(file_prediction_proto_rawDescData)
	})
	return file_prediction_proto_rawDescData
}

var file_prediction_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_prediction_proto_goTypes = []interface{}{
	(*Message)(nil),         // 0: prediction.Message
	(*MessageResponse)(nil), // 1: prediction.MessageResponse
}
var file_prediction_proto_depIdxs = []int32{
	0, // 0: prediction.Prediction.GetSubscriberPrediction:input_type -> prediction.Message
	0, // 1: prediction.Prediction.GetReverseSubscriberPrediction:input_type -> prediction.Message
	1, // 2: prediction.Prediction.GetSubscriberPrediction:output_type -> prediction.MessageResponse
	1, // 3: prediction.Prediction.GetReverseSubscriberPrediction:output_type -> prediction.MessageResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_prediction_proto_init() }
func file_prediction_proto_init() {
	if File_prediction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_prediction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_prediction_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_prediction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_prediction_proto_goTypes,
		DependencyIndexes: file_prediction_proto_depIdxs,
		MessageInfos:      file_prediction_proto_msgTypes,
	}.Build()
	File_prediction_proto = out.File
	file_prediction_proto_rawDesc = nil
	file_prediction_proto_goTypes = nil
	file_prediction_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PredictionClient is the client API for Prediction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PredictionClient interface {
	GetSubscriberPrediction(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MessageResponse, error)
	GetReverseSubscriberPrediction(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MessageResponse, error)
}

type predictionClient struct {
	cc grpc.ClientConnInterface
}

func NewPredictionClient(cc grpc.ClientConnInterface) PredictionClient {
	return &predictionClient{cc}
}

func (c *predictionClient) GetSubscriberPrediction(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/prediction.Prediction/GetSubscriberPrediction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *predictionClient) GetReverseSubscriberPrediction(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/prediction.Prediction/GetReverseSubscriberPrediction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PredictionServer is the server API for Prediction service.
type PredictionServer interface {
	GetSubscriberPrediction(context.Context, *Message) (*MessageResponse, error)
	GetReverseSubscriberPrediction(context.Context, *Message) (*MessageResponse, error)
}

// UnimplementedPredictionServer can be embedded to have forward compatible implementations.
type UnimplementedPredictionServer struct {
}

func (*UnimplementedPredictionServer) GetSubscriberPrediction(context.Context, *Message) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubscriberPrediction not implemented")
}
func (*UnimplementedPredictionServer) GetReverseSubscriberPrediction(context.Context, *Message) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReverseSubscriberPrediction not implemented")
}

func RegisterPredictionServer(s *grpc.Server, srv PredictionServer) {
	s.RegisterService(&_Prediction_serviceDesc, srv)
}

func _Prediction_GetSubscriberPrediction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PredictionServer).GetSubscriberPrediction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/prediction.Prediction/GetSubscriberPrediction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PredictionServer).GetSubscriberPrediction(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _Prediction_GetReverseSubscriberPrediction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PredictionServer).GetReverseSubscriberPrediction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/prediction.Prediction/GetReverseSubscriberPrediction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PredictionServer).GetReverseSubscriberPrediction(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _Prediction_serviceDesc = grpc.ServiceDesc{
	ServiceName: "prediction.Prediction",
	HandlerType: (*PredictionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSubscriberPrediction",
			Handler:    _Prediction_GetSubscriberPrediction_Handler,
		},
		{
			MethodName: "GetReverseSubscriberPrediction",
			Handler:    _Prediction_GetReverseSubscriberPrediction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "prediction.proto",
}

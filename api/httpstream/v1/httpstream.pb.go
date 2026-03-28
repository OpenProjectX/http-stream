// Code generated manually to match api/httpstream/v1/httpstream.proto.

package httpstreamv1

import (
	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TransferRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Source        *HTTPRequest           `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Target        *HTTPRequest           `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Pipeline      []*PipelineStage       `protobuf:"bytes,3,rep,name=pipeline,proto3" json:"pipeline,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransferRequest) Reset() {
	*x = TransferRequest{}
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransferRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*TransferRequest) ProtoMessage()    {}

func (x *TransferRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*TransferRequest) Descriptor() ([]byte, []int) {
	return file_api_httpstream_v1_httpstream_proto_rawDescGZIP(), []int{0}
}

func (x *TransferRequest) GetSource() *HTTPRequest {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *TransferRequest) GetTarget() *HTTPRequest {
	if x != nil {
		return x.Target
	}
	return nil
}

func (x *TransferRequest) GetPipeline() []*PipelineStage {
	if x != nil {
		return x.Pipeline
	}
	return nil
}

type HTTPRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Method        string                 `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Url           string                 `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Headers       map[string]string      `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ContentLength int64                  `protobuf:"varint,4,opt,name=content_length,json=contentLength,proto3" json:"content_length,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HTTPRequest) Reset() {
	*x = HTTPRequest{}
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HTTPRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*HTTPRequest) ProtoMessage()    {}

func (x *HTTPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*HTTPRequest) Descriptor() ([]byte, []int) {
	return file_api_httpstream_v1_httpstream_proto_rawDescGZIP(), []int{1}
}

func (x *HTTPRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *HTTPRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *HTTPRequest) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *HTTPRequest) GetContentLength() int64 {
	if x != nil {
		return x.ContentLength
	}
	return 0
}

type PipelineStage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Config        map[string]string      `protobuf:"bytes,2,rep,name=config,proto3" json:"config,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PipelineStage) Reset() {
	*x = PipelineStage{}
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PipelineStage) String() string { return protoimpl.X.MessageStringOf(x) }
func (*PipelineStage) ProtoMessage()    {}

func (x *PipelineStage) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*PipelineStage) Descriptor() ([]byte, []int) {
	return file_api_httpstream_v1_httpstream_proto_rawDescGZIP(), []int{2}
}

func (x *PipelineStage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PipelineStage) GetConfig() map[string]string {
	if x != nil {
		return x.Config
	}
	return nil
}

type TransferResponse struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	TransferId       string                 `protobuf:"bytes,1,opt,name=transfer_id,json=transferId,proto3" json:"transfer_id,omitempty"`
	BytesTransferred int64                  `protobuf:"varint,2,opt,name=bytes_transferred,json=bytesTransferred,proto3" json:"bytes_transferred,omitempty"`
	SourceStatusCode int32                  `protobuf:"varint,3,opt,name=source_status_code,json=sourceStatusCode,proto3" json:"source_status_code,omitempty"`
	TargetStatusCode int32                  `protobuf:"varint,4,opt,name=target_status_code,json=targetStatusCode,proto3" json:"target_status_code,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *TransferResponse) Reset() {
	*x = TransferResponse{}
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransferResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*TransferResponse) ProtoMessage()    {}

func (x *TransferResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*TransferResponse) Descriptor() ([]byte, []int) {
	return file_api_httpstream_v1_httpstream_proto_rawDescGZIP(), []int{3}
}

func (x *TransferResponse) GetTransferId() string {
	if x != nil {
		return x.TransferId
	}
	return ""
}

func (x *TransferResponse) GetBytesTransferred() int64 {
	if x != nil {
		return x.BytesTransferred
	}
	return 0
}

func (x *TransferResponse) GetSourceStatusCode() int32 {
	if x != nil {
		return x.SourceStatusCode
	}
	return 0
}

func (x *TransferResponse) GetTargetStatusCode() int32 {
	if x != nil {
		return x.TargetStatusCode
	}
	return 0
}

type HTTPRequest_HeadersEntry struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HTTPRequest_HeadersEntry) Reset() {
	*x = HTTPRequest_HeadersEntry{}
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HTTPRequest_HeadersEntry) String() string { return protoimpl.X.MessageStringOf(x) }
func (*HTTPRequest_HeadersEntry) ProtoMessage()    {}

func (x *HTTPRequest_HeadersEntry) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*HTTPRequest_HeadersEntry) Descriptor() ([]byte, []int) {
	return file_api_httpstream_v1_httpstream_proto_rawDescGZIP(), []int{1, 0}
}

type PipelineStage_ConfigEntry struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PipelineStage_ConfigEntry) Reset() {
	*x = PipelineStage_ConfigEntry{}
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PipelineStage_ConfigEntry) String() string { return protoimpl.X.MessageStringOf(x) }
func (*PipelineStage_ConfigEntry) ProtoMessage()    {}

func (x *PipelineStage_ConfigEntry) ProtoReflect() protoreflect.Message {
	mi := &file_api_httpstream_v1_httpstream_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*PipelineStage_ConfigEntry) Descriptor() ([]byte, []int) {
	return file_api_httpstream_v1_httpstream_proto_rawDescGZIP(), []int{2, 0}
}

var File_api_httpstream_v1_httpstream_proto protoreflect.FileDescriptor

var file_api_httpstream_v1_httpstream_proto_rawDescOnce sync.Once
var file_api_httpstream_v1_httpstream_proto_rawDescData = file_api_httpstream_v1_httpstream_proto_rawDesc()

func file_api_httpstream_v1_httpstream_proto_rawDescGZIP() []byte {
	file_api_httpstream_v1_httpstream_proto_rawDescOnce.Do(func() {
		file_api_httpstream_v1_httpstream_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_httpstream_v1_httpstream_proto_rawDescData)
	})
	return file_api_httpstream_v1_httpstream_proto_rawDescData
}

var file_api_httpstream_v1_httpstream_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_httpstream_v1_httpstream_proto_goTypes = []any{
	(*TransferRequest)(nil),
	(*HTTPRequest)(nil),
	(*PipelineStage)(nil),
	(*TransferResponse)(nil),
	(*HTTPRequest_HeadersEntry)(nil),
	(*PipelineStage_ConfigEntry)(nil),
}
var file_api_httpstream_v1_httpstream_proto_depIdxs = []int32{
	1,
	1,
	2,
	4,
	5,
	0,
	3,
}

func init() { file_api_httpstream_v1_httpstream_proto_init() }
func file_api_httpstream_v1_httpstream_proto_init() {
	if File_api_httpstream_v1_httpstream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_httpstream_v1_httpstream_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*TransferRequest); i {
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
		file_api_httpstream_v1_httpstream_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*HTTPRequest); i {
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
		file_api_httpstream_v1_httpstream_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PipelineStage); i {
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
		file_api_httpstream_v1_httpstream_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*TransferResponse); i {
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
		file_api_httpstream_v1_httpstream_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*HTTPRequest_HeadersEntry); i {
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
		file_api_httpstream_v1_httpstream_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*PipelineStage_ConfigEntry); i {
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

	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_httpstream_v1_httpstream_proto_rawDescData,
			NumMessages:   6,
			NumServices:   1,
		},
		GoTypes:           file_api_httpstream_v1_httpstream_proto_goTypes,
		DependencyIndexes: file_api_httpstream_v1_httpstream_proto_depIdxs,
		MessageInfos:      file_api_httpstream_v1_httpstream_proto_msgTypes,
	}.Build()

	File_api_httpstream_v1_httpstream_proto = out.File
	file_api_httpstream_v1_httpstream_proto_rawDescData = nil
	file_api_httpstream_v1_httpstream_proto_goTypes = nil
	file_api_httpstream_v1_httpstream_proto_depIdxs = nil
}

type x struct{}

func file_api_httpstream_v1_httpstream_proto_rawDesc() []byte {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    stringp("api/httpstream/v1/httpstream.proto"),
		Package: stringp("httpstream.v1"),
		Syntax:  stringp("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: stringp("github.com/OpenProjectX/http-stream/api/httpstream/v1;httpstreamv1"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: stringp("TransferRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					field("source", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.HttpRequest"),
					field("target", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.HttpRequest"),
					field("pipeline", 3, descriptorpb.FieldDescriptorProto_LABEL_REPEATED, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.PipelineStage"),
				},
			},
			{
				Name: stringp("HttpRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					field("method", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""),
					field("url", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""),
					field("headers", 3, descriptorpb.FieldDescriptorProto_LABEL_REPEATED, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.HttpRequest.HeadersEntry"),
					{
						Name:     stringp("content_length"),
						JsonName: stringp("contentLength"),
						Number:   int32p(4),
						Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
						Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_INT64),
					},
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name: stringp("HeadersEntry"),
						Field: []*descriptorpb.FieldDescriptorProto{
							field("key", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""),
							field("value", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""),
						},
						Options: &descriptorpb.MessageOptions{MapEntry: boolp(true)},
					},
				},
			},
			{
				Name: stringp("PipelineStage"),
				Field: []*descriptorpb.FieldDescriptorProto{
					field("name", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""),
					field("config", 2, descriptorpb.FieldDescriptorProto_LABEL_REPEATED, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.PipelineStage.ConfigEntry"),
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name: stringp("ConfigEntry"),
						Field: []*descriptorpb.FieldDescriptorProto{
							field("key", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""),
							field("value", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, ""),
						},
						Options: &descriptorpb.MessageOptions{MapEntry: boolp(true)},
					},
				},
			},
			{
				Name: stringp("TransferResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     stringp("transfer_id"),
						JsonName: stringp("transferId"),
						Number:   int32p(1),
						Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
						Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_STRING),
					},
					{
						Name:     stringp("bytes_transferred"),
						JsonName: stringp("bytesTransferred"),
						Number:   int32p(2),
						Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
						Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_INT64),
					},
					{
						Name:     stringp("source_status_code"),
						JsonName: stringp("sourceStatusCode"),
						Number:   int32p(3),
						Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
						Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
					},
					{
						Name:     stringp("target_status_code"),
						JsonName: stringp("targetStatusCode"),
						Number:   int32p(4),
						Label:    labelp(descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL),
						Type:     typep(descriptorpb.FieldDescriptorProto_TYPE_INT32),
					},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: stringp("StreamService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       stringp("Transfer"),
						InputType:  stringp(".httpstream.v1.TransferRequest"),
						OutputType: stringp(".httpstream.v1.TransferResponse"),
					},
				},
			},
		},
	}

	b, err := proto.Marshal(fd)
	if err != nil {
		panic(err)
	}
	return b
}

func field(name string, number int32, label descriptorpb.FieldDescriptorProto_Label, typ descriptorpb.FieldDescriptorProto_Type, typeName string) *descriptorpb.FieldDescriptorProto {
	f := &descriptorpb.FieldDescriptorProto{
		Name:   stringp(name),
		Number: int32p(number),
		Label:  labelp(label),
		Type:   typep(typ),
	}
	if typeName != "" {
		f.TypeName = stringp(typeName)
	}
	return f
}

func stringp(s string) *string { return &s }
func int32p(v int32) *int32    { return &v }
func boolp(v bool) *bool       { return &v }
func labelp(v descriptorpb.FieldDescriptorProto_Label) *descriptorpb.FieldDescriptorProto_Label {
	return &v
}
func typep(v descriptorpb.FieldDescriptorProto_Type) *descriptorpb.FieldDescriptorProto_Type {
	return &v
}

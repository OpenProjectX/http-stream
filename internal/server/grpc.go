package server

import (
	"context"
	"fmt"

	"github.com/OpenProjectX/http-stream/internal/api/httpstreamv1"
	"github.com/OpenProjectX/http-stream/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	fileDescriptor          protoreflect.FileDescriptor
	transferRequestDesc     protoreflect.MessageDescriptor
	httpRequestDesc         protoreflect.MessageDescriptor
	pipelineStageDesc       protoreflect.MessageDescriptor
	transferResponseDesc    protoreflect.MessageDescriptor
	streamServiceFullMethod = "/httpstream.v1.StreamService/Transfer"
)

func init() {
	fd, err := protodesc.NewFile(buildFileDescriptorProto(), nil)
	if err != nil {
		panic(fmt.Errorf("build protobuf descriptor: %w", err))
	}
	fileDescriptor = fd
	if err := protoregistry.GlobalFiles.RegisterFile(fd); err != nil {
		panic(fmt.Errorf("register protobuf descriptor: %w", err))
	}

	messages := fd.Messages()
	transferRequestDesc = messages.ByName("TransferRequest")
	httpRequestDesc = messages.ByName("HttpRequest")
	pipelineStageDesc = messages.ByName("PipelineStage")
	transferResponseDesc = messages.ByName("TransferResponse")
}

type GRPCServer struct {
	streamer *service.Streamer
}

func New(streamer *service.Streamer) *GRPCServer {
	return &GRPCServer{streamer: streamer}
}

func Register(grpcServer *grpc.Server, srv *GRPCServer) {
	grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: "httpstream.v1.StreamService",
		HandlerType: (*interface{})(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Transfer",
				Handler:    transferHandler,
			},
		},
		Metadata: "api/httpstream/v1/httpstream.proto",
	}, srv)
}

func transferHandler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := dynamicpb.NewMessage(transferRequestDesc)
	if err := dec(in); err != nil {
		return nil, err
	}

	handler := func(ctx context.Context, req any) (any, error) {
		internalReq, err := decodeTransferRequest(req.(*dynamicpb.Message))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		resp, err := srv.(*GRPCServer).streamer.Transfer(ctx, internalReq)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return encodeTransferResponse(resp), nil
	}

	if interceptor == nil {
		return handler(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: streamServiceFullMethod,
	}
	return interceptor(ctx, in, info, handler)
}

func decodeTransferRequest(msg *dynamicpb.Message) (*httpstreamv1.TransferRequest, error) {
	sourceMsg, err := getMessageField(msg, transferRequestDesc.Fields().ByName("source"))
	if err != nil {
		return nil, fmt.Errorf("decode source: %w", err)
	}
	targetMsg, err := getMessageField(msg, transferRequestDesc.Fields().ByName("target"))
	if err != nil {
		return nil, fmt.Errorf("decode target: %w", err)
	}

	req := &httpstreamv1.TransferRequest{
		Source: decodeHTTPRequest(sourceMsg),
		Target: decodeHTTPRequest(targetMsg),
	}

	pipelineField := transferRequestDesc.Fields().ByName("pipeline")
	list := msg.Get(pipelineField).List()
	req.Pipeline = make([]*httpstreamv1.PipelineStage, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		req.Pipeline = append(req.Pipeline, decodePipelineStage(list.Get(i).Message()))
	}

	return req, nil
}

func decodeHTTPRequest(msg protoreflect.Message) *httpstreamv1.HTTPRequest {
	fields := httpRequestDesc.Fields()
	out := &httpstreamv1.HTTPRequest{
		Method:        msg.Get(fields.ByName("method")).String(),
		URL:           msg.Get(fields.ByName("url")).String(),
		ContentLength: msg.Get(fields.ByName("content_length")).Int(),
		LocalPath:     msg.Get(fields.ByName("local_path")).String(),
	}

	headersField := fields.ByName("headers")
	headers := msg.Get(headersField).Map()
	if headers.Len() > 0 {
		out.Headers = make(map[string]string, headers.Len())
		headers.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			out.Headers[k.String()] = v.String()
			return true
		})
	}

	return out
}

func decodePipelineStage(msg protoreflect.Message) *httpstreamv1.PipelineStage {
	fields := pipelineStageDesc.Fields()
	out := &httpstreamv1.PipelineStage{
		Name: msg.Get(fields.ByName("name")).String(),
	}

	configField := fields.ByName("config")
	config := msg.Get(configField).Map()
	if config.Len() > 0 {
		out.Config = make(map[string]string, config.Len())
		config.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			out.Config[k.String()] = v.String()
			return true
		})
	}

	return out
}

func encodeTransferResponse(resp *httpstreamv1.TransferResponse) *dynamicpb.Message {
	msg := dynamicpb.NewMessage(transferResponseDesc)
	fields := transferResponseDesc.Fields()
	msg.Set(fields.ByName("transfer_id"), protoreflect.ValueOfString(resp.TransferID))
	msg.Set(fields.ByName("bytes_transferred"), protoreflect.ValueOfInt64(resp.BytesTransferred))
	msg.Set(fields.ByName("source_status_code"), protoreflect.ValueOfInt32(resp.SourceStatusCode))
	msg.Set(fields.ByName("target_status_code"), protoreflect.ValueOfInt32(resp.TargetStatusCode))
	msg.Set(fields.ByName("source_content_length"), protoreflect.ValueOfInt64(resp.SourceContentLength))
	msg.Set(fields.ByName("duration_millis"), protoreflect.ValueOfInt64(resp.DurationMillis))
	msg.Set(fields.ByName("average_bytes_per_second"), protoreflect.ValueOfFloat64(resp.AverageBytesPerSecond))
	msg.Set(fields.ByName("progress_percent"), protoreflect.ValueOfFloat64(resp.ProgressPercent))
	return msg
}

func getMessageField(msg *dynamicpb.Message, field protoreflect.FieldDescriptor) (protoreflect.Message, error) {
	if !msg.Has(field) {
		return nil, fmt.Errorf("%s is required", field.Name())
	}
	return msg.Get(field).Message(), nil
}

func buildFileDescriptorProto() *descriptorpb.FileDescriptorProto {
	return &descriptorpb.FileDescriptorProto{
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
					field("source", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.HttpRequest", "source"),
					field("target", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.HttpRequest", "target"),
					field("pipeline", 3, descriptorpb.FieldDescriptorProto_LABEL_REPEATED, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.PipelineStage", "pipeline"),
				},
			},
			{
				Name: stringp("HttpRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					field("method", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "method"),
					field("url", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "url"),
					field("headers", 3, descriptorpb.FieldDescriptorProto_LABEL_REPEATED, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.HttpRequest.HeadersEntry", "headers"),
					field("content_length", 4, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_INT64, "", "contentLength"),
					field("local_path", 5, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "localPath"),
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name: stringp("HeadersEntry"),
						Field: []*descriptorpb.FieldDescriptorProto{
							field("key", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "key"),
							field("value", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "value"),
						},
						Options: &descriptorpb.MessageOptions{MapEntry: boolp(true)},
					},
				},
			},
			{
				Name: stringp("PipelineStage"),
				Field: []*descriptorpb.FieldDescriptorProto{
					field("name", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "name"),
					field("config", 2, descriptorpb.FieldDescriptorProto_LABEL_REPEATED, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".httpstream.v1.PipelineStage.ConfigEntry", "config"),
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name: stringp("ConfigEntry"),
						Field: []*descriptorpb.FieldDescriptorProto{
							field("key", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "key"),
							field("value", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "value"),
						},
						Options: &descriptorpb.MessageOptions{MapEntry: boolp(true)},
					},
				},
			},
			{
				Name: stringp("TransferResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					field("transfer_id", 1, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_STRING, "", "transferId"),
					field("bytes_transferred", 2, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_INT64, "", "bytesTransferred"),
					field("source_status_code", 3, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_INT32, "", "sourceStatusCode"),
					field("target_status_code", 4, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_INT32, "", "targetStatusCode"),
					field("source_content_length", 5, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_INT64, "", "sourceContentLength"),
					field("duration_millis", 6, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_INT64, "", "durationMillis"),
					field("average_bytes_per_second", 7, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, "", "averageBytesPerSecond"),
					field("progress_percent", 8, descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, "", "progressPercent"),
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
}

func field(name string, number int32, label descriptorpb.FieldDescriptorProto_Label, typ descriptorpb.FieldDescriptorProto_Type, typeName, jsonName string) *descriptorpb.FieldDescriptorProto {
	f := &descriptorpb.FieldDescriptorProto{
		Name:     stringp(name),
		Number:   int32p(number),
		Label:    labelp(label),
		Type:     typep(typ),
		JsonName: stringp(jsonName),
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

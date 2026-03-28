package server

import (
	"context"

	"github.com/example/http-stream/internal/api/httpstreamv1"
	"github.com/example/http-stream/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamServiceServer interface {
	Transfer(context.Context, *httpstreamv1.TransferRequest) (*httpstreamv1.TransferResponse, error)
}

type GRPCServer struct {
	httpstreamv1Unimplemented bool
	streamer                  *service.Streamer
}

func New(streamer *service.Streamer) *GRPCServer {
	return &GRPCServer{streamer: streamer}
}

func (s *GRPCServer) Transfer(ctx context.Context, req *httpstreamv1.TransferRequest) (*httpstreamv1.TransferResponse, error) {
	resp, err := s.streamer.Transfer(ctx, req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return resp, nil
}

func Register(grpcServer *grpc.Server, srv StreamServiceServer) {
	grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: "httpstream.v1.StreamService",
		HandlerType: (*StreamServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Transfer",
				Handler: func(svc any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					in := new(httpstreamv1.TransferRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					if interceptor == nil {
						return svc.(StreamServiceServer).Transfer(ctx, in)
					}
					info := &grpc.UnaryServerInfo{
						Server:     svc,
						FullMethod: "/httpstream.v1.StreamService/Transfer",
					}
					handler := func(ctx context.Context, req any) (any, error) {
						return svc.(StreamServiceServer).Transfer(ctx, req.(*httpstreamv1.TransferRequest))
					}
					return interceptor(ctx, in, info, handler)
				},
			},
		},
		Metadata: "api/httpstream/v1/httpstream.proto",
	}, srv)
}

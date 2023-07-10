package api

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	grpcUserAgentHeader        = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *GrpcServer) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if len(md.Get(grpcGatewayUserAgentHeader)) > 0 {
			mtdt.UserAgent = md.Get(grpcGatewayUserAgentHeader)[0]
		}
		if len(md.Get(grpcUserAgentHeader)) > 0 {
			mtdt.UserAgent = md.Get(grpcUserAgentHeader)[0]
		}
		if len(md.Get(xForwardedForHeader)) > 0 {
			mtdt.ClientIP = md.Get(xForwardedForHeader)[0]
		}
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}

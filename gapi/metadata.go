package gapi

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	GRPCGatewayUserAgentHeader = "grpcgateway-user-agent"
	UserAgentHeader            = "user-agent"
	XForwardedForHeader        = "x-forwarded-for"
)

type MetaData struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *MetaData {
	mtdt := &MetaData{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("%+v", md)

		if userAgents := md.Get(GRPCGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(UserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(XForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}

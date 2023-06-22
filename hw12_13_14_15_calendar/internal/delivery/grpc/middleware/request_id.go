package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/z-sector/otus-hw/hw12_13_14_15_calendar/pkg/requestid"
)

const metadataKey = "x-request-id"

func RequestID(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "retrieving metadata is failed")
	}

	var rid string
	requestIDs := md[metadataKey]
	if len(requestIDs) > 0 {
		rid = requestIDs[len(requestIDs)-1]
	} else {
		rid = requestid.GenerateRequestID()
	}

	ctxWithReqID := requestid.CtxSetRequestID(ctx, rid)

	return handler(ctxWithReqID, req)
}

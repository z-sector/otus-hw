package requestid

import (
	"context"

	"github.com/google/uuid"
)

type ctxRequestID struct{}

func GenerateRequestID() string {
	UID, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return UID.String()
}

func CtxSetRequestID(ctx context.Context, RID string) context.Context {
	return context.WithValue(ctx, ctxRequestID{}, RID)
}

func CtxGetRequestID(ctx context.Context) string {
	if UID, ok := ctx.Value(ctxRequestID{}).(string); ok {
		return UID
	}
	return ""
}

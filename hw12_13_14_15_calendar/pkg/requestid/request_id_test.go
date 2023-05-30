package requestid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	RID := GenerateRequestID()
	require.NotEqualValues(t, "", RID)

	ctx := CtxSetRequestID(context.Background(), RID)
	actualRID := CtxGetRequestID(ctx)

	require.Equal(t, RID, actualRID)
}

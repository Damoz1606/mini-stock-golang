package gethello

import (
	"context"
	"testing"

	"github.com/Damoz1606/ministock-backend/internal/modules/query/hello/response/gethello"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name   string
	setup  func() *GetHelloHandler
	assert func(t *testing.T, result *gethello.GetHelloResponse, err error)
}

func TestGetHelloHandler_Handle(t *testing.T) {
	cases := []TestCase{
		{
			name: "Should return GetHelloResponse when Handle is called",
			setup: func() *GetHelloHandler {
				return NewGetHelloHandler()
			},
			assert: func(t *testing.T, result *gethello.GetHelloResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, "world", result.Hello)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			handler := tt.setup()
			got, err := handler.Handle(context.Background(), GetHelloQuery{})
			tt.assert(t, got, err)
		})
	}
}
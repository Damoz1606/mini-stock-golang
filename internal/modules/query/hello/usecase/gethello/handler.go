package gethello

import (
	"context"
	"github.com/Damoz1606/ministock-backend/internal/modules/query/hello/response/gethello"
)

type GetHelloHandler struct{}

func NewGetHelloHandler() *GetHelloHandler {
	return &GetHelloHandler{}
}

func (h *GetHelloHandler) Handle(ctx context.Context, query GetHelloQuery) (*gethello.GetHelloResponse, error) {
	return &gethello.GetHelloResponse{Hello: "world"}, nil
}
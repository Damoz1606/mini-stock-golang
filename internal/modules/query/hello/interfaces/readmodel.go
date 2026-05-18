package hello

import "context"

type HelloReadModel interface {
	GetHello(ctx context.Context) (string, error)
}
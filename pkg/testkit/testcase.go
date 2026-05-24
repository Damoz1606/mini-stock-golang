package testkit

import (
	"testing"
)

type TestCase[H any, R any] struct {
	Name   string
	Setup  func() H
	Assert func(t *testing.T, result R, err error)
}

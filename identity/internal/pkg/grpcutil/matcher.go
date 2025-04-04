package grpcutil

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "X-Id-Api-Key":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

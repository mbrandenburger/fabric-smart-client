package view_test

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/mock"
	view2 "github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
	"go.opentelemetry.io/otel/trace/noop"
)

var call = func(context view2.Context) (interface{}, error) {
	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

func BenchmarkView(b *testing.B) {
	//v := &views.NoopView{}
	v := &mock.View{}
	v.CallCalls(call)

	b.Run("seq", func(b *testing.B) {
		parent := &mock.ParentContext{}
		parent.ContextReturns(b.Context())
		parent.StartSpanFromReturns(b.Context(), &noop.Span{})

		for b.Loop() {
			_, _ = view.RunViewNow(parent, v, view2.WithContext(b.Context()))
		}
		reportTPS(b)
	})

	b.Run("par", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			parent := &mock.ParentContext{}
			parent.ContextReturns(b.Context())
			parent.StartSpanFromReturns(b.Context(), &noop.Span{})
			for pb.Next() {
				_, _ = view.RunViewNow(parent, v, view2.WithContext(b.Context()))
			}
		})
		reportTPS(b)
	})
}

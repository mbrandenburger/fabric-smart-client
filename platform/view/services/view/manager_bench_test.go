package view_test

import (
	"testing"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/metrics/disabled"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/mock"
	view2 "github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
	"go.opentelemetry.io/otel/trace/noop"
)

var reportTPS = func(b *testing.B) {
	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "TPS")
}

func BenchmarkManager(b *testing.B) {
	v := &mock.View{}
	v.CallCalls(call)

	registry := view.NewServiceProvider()
	idProvider := &mock.IdentityProvider{}
	idProvider.DefaultIdentityReturns([]byte("alice"))

	ch := make(chan *view2.Message)

	session := &mock.Session{}
	session.ReceiveReturns(ch)

	commLayer := mock.CommLayer{}
	commLayer.MasterSessionReturns(session, nil)

	vm := view.NewManager(registry, &commLayer, &mock.EndpointService{}, idProvider, view.NewRegistry(), noop.NewTracerProvider(), &disabled.Provider{}, nil)

	b.Run("seq", func(b *testing.B) {
		for b.Loop() {
			_, _ = vm.InitiateView(v, b.Context())
		}
		reportTPS(b)
	})

	b.Run("par", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = vm.InitiateView(v, b.Context())
			}
		})
		reportTPS(b)
	})
}

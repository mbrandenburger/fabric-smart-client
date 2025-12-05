package test

import (
	"testing"
	"time"

	"github.com/hyperledger-labs/fabric-smart-client/integration/benchmark"
	"github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple/views"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/grpc"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/metrics/disabled"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/grpc/client"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/grpc/server"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/grpc/server/protos"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

func Benchmark(b *testing.B) {
	setupServer(b)

	b.Run("seq", func(b *testing.B) {
		cli := setupClient(b)

		for b.Loop() {
			resp, err := cli.CallView("fid", nil)
			require.NoError(b, err)
			require.NotNil(b, resp)
		}
		b.ReportAllocs()
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "TPS")
	})

	b.Run("parallel-1", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			cli := setupClient(b)

			for pb.Next() {
				resp, err := cli.CallView("fid", nil)
				require.NoError(b, err)
				require.NotNil(b, resp)
			}
		})
		b.ReportAllocs()
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "TPS")
	})
}

type caller interface {
	CallView(string, []byte) (any, error)
	CreateSignedCommand(payload interface{}, signingIdentity client.SigningIdentity) (*protos.SignedCommand, error)
}

func setupServer(tb testing.TB) {
	tb.Helper()

	mDefaultIdentity := view.Identity("server identity")
	mSigner := &benchmark.MockSigner{
		SerializeFunc: func() ([]byte, error) {
			return mDefaultIdentity.Bytes(), nil
		},
		SignFunc: func(bytes []byte) ([]byte, error) {
			return bytes, nil
		}}
	mIdentityProvider := &benchmark.MockIdentityProvider{DefaultSigner: mDefaultIdentity}
	mSigService := &benchmark.MockSignerProvider{DefaultSigner: mSigner}

	// marshaller
	tm, err := server.NewResponseMarshaler(mIdentityProvider, mSigService)
	require.NoError(tb, err)
	require.NotNil(tb, tm)

	// setup server
	grpcSrv, err := grpc.NewGRPCServer("localhost:8080", grpc.ServerConfig{
		ConnectionTimeout: 0,
		SecOpts: grpc.SecureOptions{
			Certificate: certPEM,
			Key:         keyPEM,
			UseTLS:      true,
		},
		KaOpts:             grpc.KeepaliveOptions{},
		Logger:             nil,
		HealthCheckEnabled: false,
	})

	require.NoError(tb, err)
	require.NotNil(tb, grpcSrv)

	tb.Logf("listening on %v", grpcSrv.Listener().Addr())

	srv, err := server.NewViewServiceServer(tm, &server.YesPolicyChecker{}, server.NewMetrics(&disabled.Provider{}), noop.NewTracerProvider())
	require.NoError(tb, err)
	require.NotNil(tb, srv)

	// our view manager
	vm := &benchmark.MockViewManager{Constructor: func() view.View {
		return &views.SimpleView{}
	}}

	// register view manager wit grpc impl
	server.InstallViewHandler(vm, srv, noop.NewTracerProvider())

	// register grpc impl with grpc server
	protos.RegisterViewServiceServer(grpcSrv.Server(), srv)

	// start the actual grpc server
	go func() {
		_ = grpcSrv.Start()
	}()
	tb.Cleanup(grpcSrv.Stop)

	return
}

func setupClient(tb testing.TB) caller {
	tb.Helper()

	mDefaultIdentity := view.Identity("client identity")
	mSigner := &benchmark.MockSigner{
		SerializeFunc: func() ([]byte, error) {
			return mDefaultIdentity.Bytes(), nil
		},
		SignFunc: func(bytes []byte) ([]byte, error) {
			return bytes, nil
		}}

	// setup client
	cfg := &client.Config{
		ID: "someID",
		ConnectionConfig: &grpc.ConnectionConfig{
			Address:            "localhost:8080",
			ConnectionTimeout:  0,
			TLSEnabled:         false,
			TLSClientSideAuth:  false,
			TLSDisabled:        true,
			TLSRootCertFile:    "",
			TLSRootCertBytes:   nil,
			ServerNameOverride: "",
			Usage:              "",
		},
	}

	grpcClient, err := grpc.NewGRPCClient(grpc.ClientConfig{
		SecOpts: grpc.SecureOptions{
			ServerRootCAs: [][]byte{certPEM},
			UseTLS:        true,
		},
		KaOpts:       grpc.KeepaliveOptions{},
		Timeout:      5 * time.Second,
		AsyncConnect: false,
	})
	require.NoError(tb, err)
	require.NotNil(tb, grpcClient)

	conn, err := grpcClient.NewConnection("localhost:8080")
	require.NoError(tb, err)
	require.NotNil(tb, conn)

	cli, err := client.NewClient2(
		grpcClient,
		cfg,
		mSigner,
		noop.NewTracerProvider(),
	)
	require.NoError(tb, err)
	require.NotNil(tb, cli)

	return cli
}

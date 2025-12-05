package test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"

	simpleviews "github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple/views"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/grpc/server/protos"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
	"github.com/hyperledger/fabric-lib-go/common/flogging"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	keyPEM  []byte
	certPEM []byte
)

func init() {
	flogging.Init(flogging.Config{
		LogSpec: "error",
	})

	keyPEM, certPEM, _ = makeSelfSignedCert()
}

type serverImpl struct {
	protos.UnimplementedViewServiceServer

	workload view.View
}

func (s *serverImpl) ProcessCommand(ctx context.Context, command *protos.SignedCommand) (*protos.SignedCommandResponse, error) {
	resp, err := s.workload.Call(nil)
	if err != nil {
		return nil, err
	}

	// TODO include resp in signed response
	_ = resp

	return &protos.SignedCommandResponse{}, nil
}

func (s *serverImpl) StreamCommand(g grpc.BidiStreamingServer[protos.SignedCommand, protos.SignedCommandResponse]) error {
	//TODO implement me
	panic("implement me")
}

func BenchmarkReferenceImpl(b *testing.B) {
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		b.Fatalf("failt to create x509 keypair: %v", err)
	}

	serverTLS := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		// Require and verify client certs if you want mTLS:
		// ClientAuth: tls.RequireAndVerifyClientCert,
	}

	// setup server
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		b.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(credentials.NewTLS(serverTLS)))
	grpcServer := grpc.NewServer(opts...)

	srv := &serverImpl{workload: &simpleviews.SimpleView{}}

	protos.RegisterViewServiceServer(grpcServer, srv)
	go grpcServer.Serve(lis)

	// prepare client-side root pool trusting the self-signed cert
	rootPool := x509.NewCertPool()
	ok := rootPool.AppendCertsFromPEM(certPEM)
	if !ok {
		b.Fatalf("failed to append cert to pool")
	}

	clientTLS := &tls.Config{
		RootCAs:    rootPool,
		ServerName: "localhost", // must match cert's DNSNames / SAN
	}

	// and client
	var clientOpts []grpc.DialOption
	clientOpts = append(clientOpts, grpc.WithTransportCredentials(credentials.NewTLS(clientTLS)))
	conn, err := grpc.NewClient("localhost:8080", clientOpts...)
	if err != nil {
		b.Fatalf("failed to create client conn: %v", err)

	}
	defer conn.Close()

	client := protos.NewViewServiceClient(conn)

	b.Run("seq", func(b *testing.B) {
		for b.Loop() {
			resp, err := client.ProcessCommand(b.Context(), &protos.SignedCommand{
				Command:   nil,
				Signature: nil,
			})
			require.NoError(b, err)
			require.NotNil(b, resp)
		}
		b.ReportAllocs()
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "TPS")
	})

	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				resp, err := client.ProcessCommand(b.Context(), &protos.SignedCommand{
					Command:   nil,
					Signature: nil,
				})
				require.NoError(b, err)
				require.NotNil(b, resp)
			}
		})
		b.ReportAllocs()
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "TPS")
	})

	b.Run("parallel-mc", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			var clientOpts []grpc.DialOption
			clientOpts = append(clientOpts, grpc.WithTransportCredentials(credentials.NewTLS(clientTLS)))
			conn, err := grpc.NewClient("localhost:8080", clientOpts...)
			if err != nil {
				b.Fatalf("failed to create client conn: %v", err)

			}
			defer conn.Close()

			client := protos.NewViewServiceClient(conn)

			for pb.Next() {
				resp, err := client.ProcessCommand(b.Context(), &protos.SignedCommand{
					Command:   nil,
					Signature: nil,
				})
				require.NoError(b, err)
				require.NotNil(b, resp)
			}
		})
		b.ReportAllocs()
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "TPS")
	})

}

// makeSelfSignedCert generates a localhost self-signed cert using ECDSA P-256.
// It returns the tls.Certificate and the PEM-encoded cert for the client root pool.
func makeSelfSignedCert() ([]byte, []byte, error) {
	// 1. generate ECDSA private key
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// 2. certificate template
	serial, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Local Test CA"},
		},
		NotBefore:             time.Now().Add(-time.Minute),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"localhost"},
	}

	// 3. sign it (self-signed)
	derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	// 4. PEM-encode cert & private key
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	return keyPEM, certPEM, nil
}

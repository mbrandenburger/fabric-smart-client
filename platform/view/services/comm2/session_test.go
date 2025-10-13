/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package comm2

import (
	"context"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/hyperledger-labs/fabric-smart-client/platform/common/services/logging"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const (
	timeout = 5 * time.Second
	tick    = 200 * time.Millisecond
)

func TestSession(t *testing.T) {
	// let check that at the end of this test all our go routines are stopped
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	logging.Init(logging.Config{
		LogSpec: "grpc=error:debug",
	})

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener for test server: %v", err)
	}
	address := lis.Addr().String()

	ctx, cancel := context.WithCancel(t.Context())
	openingListener, serverDone := setupServer(t, ctx, lis)

	cm := setupClient(t)
	defer cm.Close()

	cs := setupService(t, ctx, cm, address)
	clientSession, err := cs.NewSession("mario", "comeCtx", "", nil)
	require.NoError(t, err)

	var serverSession *Session
	require.EventuallyWithT(t, func(tc *assert.CollectT) {
		serverSession = <-openingListener
		assert.NotNil(tc, serverSession)
	}, timeout, tick)

	msg := []byte("hello")

	var wg sync.WaitGroup
	wg.Add(100)
	for range 100 {
		go func() {
			defer wg.Done()
			for ctx.Err() == nil {
				err := clientSession.SendWithContext(ctx, msg)
				if err != nil {
					assert.ErrorIs(t, err, context.Canceled)
				}
			}
		}()
	}

	wg.Add(90)
	for range 90 {
		go func() {
			defer wg.Done()
			for ctx.Err() == nil {
				select {
				case <-ctx.Done():
				case m := <-serverSession.Receive():
					assert.Equal(t, msg, m.Payload)
				}
			}
		}()
	}

	// let's give the producer a bit time
	runtime.Gosched()
	for {
		value := rand.Intn(100000)
		if value == 0 {
			break
		}
	}

	logger.Warnf("time to close")

	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
			clientSession.Close()
		}()
	}

	cancel()
	<-serverDone
	wg.Wait()
}

func setupClient(t testing.TB) *grpc.Client {
	t.Helper()

	c := grpc.ClientConfig{
		Timeout:      timeout,
		AsyncConnect: true,
	}

	cm, err := grpc.NewGRPCClient(c)
	require.NoError(t, err)

	return cm
}

func setupService(t testing.TB, ctx context.Context, cm *grpc.Client, address string) *Service {
	t.Helper()

	cc, err := cm.NewConnection(address)
	require.NoError(t, err)

	s := &Service{
		clientLookup: func(endpoint string, pkid []byte) (P2PServiceClient, error) {
			return NewP2PServiceClient(cc), nil
		},
		ctx: ctx,
	}

	return s
}

func setupServer(t testing.TB, ctx context.Context, lis net.Listener) (chan *Session, chan struct{}) {
	t.Helper()

	c := grpc.ServerConfig{SecOpts: grpc.SecureOptions{UseTLS: false}}

	srv, err := grpc.NewGRPCServerFromListener(lis, c)
	require.NoError(t, err)

	openingListener := make(chan *Session)

	p2pServer := &server{
		openingListener: openingListener,
	}
	RegisterP2PServiceServer(srv.Server(), p2pServer)

	done := make(chan struct{})

	go func() {
		_ = srv.Start()
	}()

	go func() {
		<-ctx.Done()
		srv.Stop()
		close(done)
	}()

	return openingListener, done
}

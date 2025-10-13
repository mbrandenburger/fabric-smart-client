/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package comm2

import (
	"context"
	"net"
	sync "sync"
	"testing"

	"github.com/hyperledger-labs/fabric-smart-client/platform/common/services/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkSession(b *testing.B) {
	logging.Init(logging.Config{
		LogSpec: "grpc=error:error",
	})

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("failed to create listener for test server: %v", err)
	}
	address := lis.Addr().String()

	ctx, cancel := context.WithCancel(b.Context())
	openingListener, serverDone := setupServer(b, ctx, lis)

	cm := setupClient(b)
	defer cm.Close()

	cs := setupService(b, ctx, cm, address)
	clientSession, err := cs.NewSession("mario", "someCtx", "", nil)
	require.NoError(b, err)

	var serverSession *Session
	require.EventuallyWithT(b, func(tc *assert.CollectT) {
		serverSession = <-openingListener
		assert.NotNil(tc, serverSession)
	}, timeout, tick)

	logger.Warnf("getting ready")

	msg := []byte("hello")
	msg2 := []byte("moin")
	var wg sync.WaitGroup
	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
			for ctx.Err() == nil {
				select {
				case <-ctx.Done():
				case m := <-serverSession.Receive():
					assert.Equal(b, msg, m.Payload)
					assert.NoError(b, serverSession.SendWithContext(ctx, msg2))
				}
			}
		}()
	}

	require.NoError(b, clientSession.SendWithContext(ctx, msg))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = clientSession.SendWithContext(ctx, msg)
		if err != nil {
			assert.ErrorIs(b, err, context.Canceled)
		}
		m := <-clientSession.Receive()
		assert.Equal(b, msg2, m.Payload)
	}
	b.StopTimer()

	clientSession.Close()
	cancel()
	<-serverDone
	wg.Wait()
}

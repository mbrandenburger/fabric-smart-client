/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm2

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type server struct {
	UnimplementedP2PServiceServer

	openingListener chan *Session
}

func (s *server) OpenSessionStream(stream P2PService_OpenSessionStreamServer) error {
	logger.Debugf("New OpenSessionStream: %v", stream)

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok || len(md["caller"]) == 0 {
		return errors.New("no caller defined in metadata")
	}

	callerViewID := md.Get("caller")[0]
	contextID := md.Get("contextID")[0]
	sessionID := md.Get("sessionID")[0]
	// TODO: validate endpoint and endpointPKID
	fromEndpoint := md.Get("endpoint")[0]
	fromEndpointPKID, err := base64.StdEncoding.DecodeString(md.Get("endpointPKID")[0])
	if err != nil {
		return err
	}

	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return fmt.Errorf("failed loading peer info from stream context")
	}

	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if ok {
		// TODO: check endpoint end endpontPKID against TLS certs
		logger.Debugf("grpc client: %v - %v", p, tlsInfo.AuthType())
	} else {
		logger.Debugf("no TLS! be careful")
	}

	// this is our fresh session
	ss := &Session{
		sessionID:        sessionID,
		contextID:        contextID,
		callerViewID:     callerViewID,
		fromEndpoint:     fromEndpoint,
		fromEndpointPKID: fromEndpointPKID,
		//
		sendCh:    make(chan *ViewPacket),
		receiveCh: make(chan *view.Message),
		closing:   make(chan struct{}, 1),
		closed:    make(chan struct{}),
	}

	// connect our listener
	go func() {
		select {
		case <-stream.Context().Done():
		case s.openingListener <- ss:
			logger.Debugf("new session (server-side)")
		}
	}()

	// we keep this call open until the session is closed; otherwise we terminate the grpc
	return ss.run(stream)
}

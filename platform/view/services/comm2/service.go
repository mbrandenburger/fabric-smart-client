/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm2

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger-labs/fabric-smart-client/pkg/utils"
	"github.com/hyperledger-labs/fabric-smart-client/pkg/utils/errors"
	"github.com/hyperledger-labs/fabric-smart-client/platform/common/services/logging"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/endpoint"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/grpc"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	masterSession = "master of puppets I'm pulling your strings"
)

var logger = logging.MustGetLogger()

type Service struct {
	clientLookup func(endpoint string, pkid []byte) (P2PServiceClient, error)

	// we register our grpc service
	*server

	baseCtx context.Context

	masterSession view.Session

	sessions    sync.Map
	connections sync.Map

	certCache sync.Map
}

func NewService(ep *endpoint.Service) *Service {
	// setup our master session
	ms := &Session{
		sessionID: masterSession,

		sendCh:    make(chan *ViewPacket),
		receiveCh: make(chan *view.Message),
		closing:   make(chan struct{}, 1),
		closed:    make(chan struct{}),
	}

	ctx := context.TODO()
	openingListener := make(chan *Session, 1)

	s := &Service{
		server: &server{
			openingListener: openingListener,
		},
		baseCtx:       ctx,
		masterSession: ms,
	}

	s.clientLookup = func(endpointAddr string, pkid []byte) (P2PServiceClient, error) {
		logger.Debugf("client Lookup: %v %x", endpointAddr, pkid)

		var (
			client P2PServiceClient
			err    error
		)

		if len(endpointAddr) == 0 {
			return nil, errors.New("endpoint empty")
		}

		// check if we already have a connection
		cl, exists := s.connections.Load(endpointAddr)
		if exists {
			return cl.(P2PServiceClient), nil
		}

		// if not create one ...
		var certs [][]byte
		for _, r := range ep.Resolvers() {
			if len(r.TLSRootCa) != 0 {
				certs = append(certs, r.TLSRootCa)
			}
		}

		c := grpc.ClientConfig{
			SecOpts: grpc.SecureOptions{
				ServerRootCAs: certs,
				UseTLS:        true,
			},
			Timeout:      10 * time.Second,
			AsyncConnect: false,
		}

		var cm *grpc.Client
		cm, err = grpc.NewGRPCClient(c)
		if err != nil {
			return nil, err
		}

		var cc *grpc2.ClientConn
		cc, err = cm.NewConnection(endpointAddr)
		if err != nil {
			return nil, err
		}

		client = NewP2PServiceClient(cc)

		actual, loaded := s.connections.LoadOrStore(endpointAddr, client)
		if loaded {
			cm.Close()
			return actual.(P2PServiceClient), nil
		}

		return client, err
	}

	// let's connect the opening Listener with the master session
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Errorf("done: %v", ctx.Err())
				return
			case l := <-openingListener:
				// TODO: what to do with the error here?
				_ = safeSession(&s.sessions, l.sessionID, l)
				// read first message
				// and forward to master session
				msg := <-l.Receive()
				ms.receiveCh <- msg
				logger.Debugf("forwarded msg to master: %v", msg)
			}
		}
	}()

	return s
}

func (s *Service) NewSessionWithID(sessionID, contextID, endpoint string, pkid []byte, caller view.Identity, msg *view.Message) (view.Session, error) {
	sess, exists := checkSessionExists(&s.sessions, sessionID)
	if !exists {
		// So far this call is used when a session was established via the network and re are
		// that means that the session must exist already
		panic("programming error")
	}

	// TODO: this is so stupid ... and may cause issues - please fix me
	go func() {
		ss := sess.(*Session)
		// TODO: DON'T DO THIS!!!! how can we do this better here?
		ss.caller = caller
		if msg != nil {
			ss.receiveCh <- msg
		}
	}()

	return sess, nil
}

func (s *Service) base64PkId(bytes []byte) string {
	key := string(bytes)
	if v, ok := s.certCache.Load(key); ok {
		return v.(string)
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	s.certCache.Store(key, encoded)
	return encoded
}

func (s *Service) NewSession(callerViewID string, contextID string, endpoint string, pkid []byte) (view.Session, error) {
	sessionID := utils.GenerateUUID()

	sess, exists := checkSessionExists(&s.sessions, sessionID)
	if exists {
		return sess, nil
	}

	logger.Debugf("New Session: %v %v %v %x", callerViewID, contextID, endpoint, pkid)
	if s.clientLookup == nil {
		panic("programming error: service not correctly set up")
	}

	c, err := s.clientLookup(endpoint, pkid)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(s.baseCtx)

	// we set the session metadata
	ctx = metadata.AppendToOutgoingContext(ctx,
		"caller", callerViewID,
		"contextID", contextID,
		"sessionID", sessionID,
		"endpoint", endpoint,
		//"endpointPKID", base64.StdEncoding.EncodeToString(pkid),
		"endpointPKID", s.base64PkId(pkid),
	)

	sc, err := c.OpenSessionStream(ctx)
	if err != nil {
		return nil, err
	}

	ss := &Session{
		// this is our session metadata
		sessionID:    sessionID,
		contextID:    contextID,
		callerViewID: callerViewID,
		//
		sendCh:    make(chan *ViewPacket),
		receiveCh: make(chan *view.Message),
		closing:   make(chan struct{}, 1),
		closed:    make(chan struct{}),
		//
		closerF: cancel,
	}
	logger.Debugf("new session (client-side)")

	if err := safeSession(&s.sessions, sessionID, ss); err != nil {
		return nil, err
	}

	go func() {
		err := ss.run(sc)
		logger.Debugf("client session done: %v", err)
	}()
	return ss, err
}

func (s *Service) MasterSession() (view.Session, error) {
	// TODO: can't we just remove the concept of a master session?
	return s.masterSession, nil
}

func (s *Service) DeleteSessions(ctx context.Context, sessionID string) {
	_ = deleteSession(&s.sessions, sessionID)
}

func checkSessionExists(sessions *sync.Map, sessionID string) (view.Session, bool) {
	sess, exists := sessions.Load(sessionID)
	if !exists {
		return nil, exists
	}
	return sess.(*Session), exists
}

func safeSession(sessions *sync.Map, sessionID string, sess *Session) error {
	if _, loaded := sessions.LoadOrStore(sessionID, sess); loaded {
		return fmt.Errorf("session for %v already exists", sessionID)
	}
	return nil
}

func deleteSession(sessions *sync.Map, sessionID string) error {
	se, loaded := sessions.LoadAndDelete(sessionID)
	if !loaded {
		// ignore ...
		return nil
	}
	sess := se.(*Session)
	sess.Close()
	return nil
}

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package comm2

import (
	"context"

	"github.com/hyperledger-labs/fabric-smart-client/pkg/utils/errors"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
	"golang.org/x/sync/errgroup"
)

type Session struct {
	sessionID        string
	contextID        string
	caller           view.Identity
	callerViewID     string
	fromEndpoint     string
	fromEndpointPKID []byte

	sendCh    chan *ViewPacket
	receiveCh chan *view.Message

	closing chan struct{}
	closed  chan struct{}
}

func (s *Session) Info() view.SessionInfo {
	return view.SessionInfo{
		ID:           s.sessionID,
		Caller:       s.caller,
		CallerViewID: s.callerViewID,
		Endpoint:     s.fromEndpoint,
		EndpointPKID: s.fromEndpointPKID,
		Closed:       s.isClosed(),
	}
}

func (s *Session) Send(payload []byte) error {
	return s.SendWithContext(context.TODO(), payload)
}

func (s *Session) SendWithContext(ctx context.Context, payload []byte) error {
	return s.send(ctx, &ViewPacket{Status: view.OK, Payload: payload})
}

func (s *Session) SendError(payload []byte) error {
	return s.SendErrorWithContext(context.TODO(), payload)
}

func (s *Session) SendErrorWithContext(ctx context.Context, payload []byte) error {
	return s.send(ctx, &ViewPacket{Status: view.ERROR, Payload: payload})
}

func (s *Session) Receive() <-chan *view.Message {
	return s.receiveCh
}

func (s *Session) Close() {
	select {
	case s.closing <- struct{}{}:
		logger.Debugf("closing session")
		close(s.closed)
		<-s.closed
	case <-s.closed:
	}
	logger.Debugf("session closed")
}

func (s *Session) isClosed() bool {
	select {
	case <-s.closed:
		return true
	default:
	}

	return false
}

func (s *Session) send(ctx context.Context, p *ViewPacket) error {
	select {
	case <-s.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	select {
	case <-s.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case s.sendCh <- p:
		return nil
	}
}

func (s *Session) run(stream grpcStream) error {
	g, gCtx := errgroup.WithContext(stream.Context())

	// run the receiver
	g.Go(func() error {
		defer s.Close()
		for gCtx.Err() == nil {
			req, err := stream.Recv()
			if err != nil {
				return errors.Wrap(err, "error receiving request")
			}

			msg := &view.Message{
				SessionID:    s.sessionID,
				ContextID:    s.contextID,
				Caller:       s.callerViewID, // TODO: should we rename the caller in the view.Message?
				FromEndpoint: s.fromEndpoint,
				FromPKID:     s.fromEndpointPKID,
				Status:       req.Status,
				Payload:      req.Payload,
				Ctx:          nil,
			}

			select {
			case <-gCtx.Done():
				// Do nothing here as next loop iteration will gCtx.Err() be non-nil.
			case s.receiveCh <- msg:
				logger.Debugf("Received msg: %v", msg)
			}
		}
		return gCtx.Err()
	})

	// run the sender
	g.Go(func() error {
		defer s.Close()
		for gCtx.Err() == nil {
			select {
			case <-gCtx.Done():
				// Do nothing here as next loop iteration will gCtx.Err() be non-nil.
			case res, ok := <-s.sendCh:
				if !ok {
					break
				}

				if err := stream.Send(res); err != nil {
					return err
				}
				logger.Debugf("Sent msg: %v", res)
			}
		}
		return gCtx.Err()
	})

	return g.Wait()
}

type grpcStream interface {
	Send(*ViewPacket) error
	Recv() (*ViewPacket, error)
	Context() context.Context
}

package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

type FoxParams struct {
}

type FoxView struct {
	params FoxParams
}

func (q *FoxView) Call(viewCtx view.Context) (interface{}, error) {
	s := viewCtx.Session()

	timeout := time.After(5 * time.Second)
	select {
	case msg := <-s.Receive():
		// a read from ch has occurred
		_ = msg

	case <-timeout:
		// the read from ch has timed out
		return nil, errors.New("receive timeout, disappointing")
	}

	err := s.Send([]byte("fox got you"))
	if err != nil {
		return nil, fmt.Errorf("error sending response: %w", err)
	}

	return "OK", nil
}

type FoxViewFactory struct{}

func (c *FoxViewFactory) NewView(in []byte) (view.View, error) {
	f := &FoxView{}
	if err := json.Unmarshal(in, &f.params); err != nil {
		return nil, err
	}

	return f, nil
}

package views

import (
	"encoding/json"
	"time"

	"github.com/hyperledger-labs/fabric-smart-client/pkg/utils/errors"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/id"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

type SquirrelParams struct {
}

type SquirrelView struct {
	params SquirrelParams
}

func (q *SquirrelView) Call(viewCtx view.Context) (interface{}, error) {

	identityProvider, err := id.GetProvider(viewCtx)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting identity provider")
	}

	responder := identityProvider.Identity("simple")

	s, err := viewCtx.GetSession(viewCtx.Initiator(), responder)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting session")
	}
	//defer s.Close()

	err = s.Send([]byte("Hello"))
	if err != nil {
		return nil, errors.Wrapf(err, "error sending message to responder")
	}

	timeout := time.After(5 * time.Second)
	select {
	case msg := <-s.Receive():
		// a read from ch has occurred
		_ = msg
		//logger.Warnf("recv: %s", msg.Payload)
	case <-timeout:
		// the read from ch has timed out
		return nil, errors.New("receive timeout, disappointing")
	}

	return "OK", nil
}

type SquirrelViewFactory struct{}

func (c *SquirrelViewFactory) NewView(in []byte) (view.View, error) {
	f := &SquirrelView{}
	if err := json.Unmarshal(in, &f.params); err != nil {
		return nil, err
	}

	return f, nil
}

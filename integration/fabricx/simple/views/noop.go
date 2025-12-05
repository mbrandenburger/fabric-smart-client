/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package views

import (
	"encoding/json"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

type NoopParams struct {
}

type NoopView struct {
	params NoopParams
}

func (q *NoopView) Call(viewCtx view.Context) (interface{}, error) {
	return "OK", nil
}

type NoopViewFactory struct{}

func (c *NoopViewFactory) NewView(in []byte) (view.View, error) {
	f := &NoopView{}
	if err := json.Unmarshal(in, &f.params); err != nil {
		return nil, err
	}

	return f, nil
}

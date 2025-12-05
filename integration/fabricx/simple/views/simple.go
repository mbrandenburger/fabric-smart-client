package views

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

func init() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	pr = privateKey

}

var pr *ecdsa.PrivateKey
var hash [32]byte

type SimpleParams struct {
}

type SimpleView struct {
	params NoopParams

	skipSignature bool
}

func (q *SimpleView) Call(viewCtx view.Context) (interface{}, error) {
	msg := "hello, world"
	hash = sha256.Sum256([]byte(msg))

	// we can run this workload without the signature verification
	if q.skipSignature {
		return base64.StdEncoding.EncodeToString(hash[:]), nil
	}

	sig, err := ecdsa.SignASN1(rand.Reader, pr, hash[:])
	if err != nil {
		return "error", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

type SimpleViewFactory struct{}

func (c *SimpleViewFactory) NewView(in []byte) (view.View, error) {
	f := &SimpleView{}
	if err := json.Unmarshal(in, &f.params); err != nil {
		return nil, err
	}

	return f, nil
}

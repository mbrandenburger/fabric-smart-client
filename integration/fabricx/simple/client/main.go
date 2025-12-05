package main

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"os"
	"path"
	"time"

	"github.com/hyperledger-labs/fabric-smart-client/integration/benchmark"
	simpleviews "github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple/views"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/grpc"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/grpc/client"
	view2 "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/grpc/client/cmd"
	"go.opentelemetry.io/otel/trace/noop"
)

func main() {

	cli, err := setupClient("../app/out/testdata/fsc/nodes/perf.0/client-config.yaml")
	if err != nil {
		log.Fatalf("error setting up the client %v", err)
	}

	cfg := benchmark.Config{
		NumWorker:       10,
		Duration:        30 * time.Second,
		WarumupDuration: 5 * time.Second,
		MaxSleep:        0,
	}

	benchmark.Benchmark(cfg, func() {
		//callSquirrel(cli)
		callSimple(cli)
		//callNoop(cli)
	})
}

func callNoop(cli caller) {
	params := &simpleviews.NoopParams{}
	input, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("error marshalling view input %v", err)
	}

	execute(cli, "noop", input)
}

func callSimple(cli caller) {
	params := &simpleviews.SimpleParams{}
	input, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("error marshalling view input %v", err)
	}

	execute(cli, "simple", input)
}

func callSquirrel(cli caller) {
	params := &simpleviews.SquirrelParams{}
	input, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("error marshalling view input %v", err)
	}

	execute(cli, "squirrel", input)
}

func execute(cli caller, viewName string, input []byte) {
	resp, err := cli.CallView(viewName, input)
	if err != nil {
		log.Fatalf("error calling view %v", err)
	}

	var result string
	err = json.Unmarshal(resp.([]byte), &result)
	if err != nil {
		log.Fatalf("error unmarshalling result %v", err)
	}

	_ = result
	//log.Printf("Result: %v\n", result)
}

type caller interface {
	CallView(string, []byte) (interface{}, error)
}

func setupClient(configFile string) (caller, error) {
	config, err := view2.ConfigFromFile(configFile)
	if err != nil {
		return nil, err
	}

	cert, err := os.ReadFile(config.SignerConfig.IdentityPath)

	signer := &benchmark.MockSigner{
		SignFunc: func(bytes []byte) ([]byte, error) {
			h := sha256.Sum256(bytes)
			return h[:], nil
		},
		SerializeFunc: func() ([]byte, error) {
			return cert, nil
		}}

	cc := &grpc.ConnectionConfig{
		Address:           config.Address,
		TLSEnabled:        true,
		TLSRootCertFile:   path.Join(config.TLSConfig.PeerCACertPath),
		ConnectionTimeout: 10 * time.Second,
	}

	c, err := client.NewClient(
		&client.Config{
			ConnectionConfig: cc,
		},
		signer,
		noop.NewTracerProvider(),
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

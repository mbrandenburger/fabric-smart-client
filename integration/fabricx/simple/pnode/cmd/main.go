/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple"
	simpleviews "github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple/views"
	"github.com/hyperledger-labs/fabric-smart-client/node"
	viewregistry "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"
)

func main() {
	cwd, _ := os.Getwd()
	pth := flag.String("conf", cwd, "the directory that contains the core.yaml configuration file")
	flag.Parse()

	fsc, err := startFSC(*pth, path.Join(*pth, "data"))
	if err != nil {
		log.Fatal(err)
	}

	// Register views and responders (communication with other FSC nodes)
	reg := viewregistry.GetRegistry(fsc)
	reg.RegisterResponder(&simpleviews.FoxView{}, &simpleviews.SquirrelView{})

	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	// Stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fsc.Stop()
}

func startFSC(confPath, datadir string) (*node.Node, error) {
	if len(datadir) != 0 {
		if err := os.MkdirAll(datadir, 0755); err != nil {
			return nil, fmt.Errorf("error creating data directory %s: %w", datadir, err)
		}
	}

	fsc := node.NewWithConfPath(confPath)
	if err := fsc.InstallSDK(simple.NewSDK(fsc)); err != nil {
		return nil, fmt.Errorf("error installing fsc: %w", err)
	}
	fmt.Printf("start the start ...\n")
	if err := fsc.Start(); err != nil {
		return nil, fmt.Errorf("error starting fsc: %w", err)
	}
	fmt.Printf("Started peer with ID=[%s]\n", fsc.ID())

	return fsc, nil
}

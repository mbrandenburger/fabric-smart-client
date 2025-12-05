/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/hyperledger-labs/fabric-smart-client/integration/benchmark"
	"github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple"
	simpleviews "github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple/views"
	"github.com/hyperledger-labs/fabric-smart-client/node"
	viewregistry "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view"
)

func main() {
	cwd, _ := os.Getwd()
	pth := flag.String("conf", cwd, "the directory that contains the core.yaml configuration file")
	port := flag.String("port", "9000", "the API port for the application")
	flag.Parse()

	fsc, err := startFSC(*pth, path.Join(*pth, "data"))
	if err != nil {
		log.Fatal(err)
	}

	// Register views and responders (communication with other FSC nodes)
	reg := viewregistry.GetRegistry(fsc)
	reg.RegisterFactory("simple", &simpleviews.SimpleViewFactory{})
	reg.RegisterFactory("squirrel", &simpleviews.SquirrelViewFactory{})
	reg.RegisterFactory("noop", &simpleviews.NoopViewFactory{})
	reg.RegisterResponder(&simpleviews.FoxView{}, &simpleviews.SquirrelView{})

	// get the view manager
	vm, err := viewregistry.GetManager(fsc)
	if err != nil {
		log.Fatal(err)
	}

	// register invoker hook
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{n}", func(w http.ResponseWriter, r *http.Request) {
		nStr := r.PathValue("n")

		n, err := strconv.Atoi(nStr)
		if err != nil {
			// Invalid integer → return 400 Bad Request
			http.Error(w, "n must be an integer", http.StatusBadRequest)
			return
		}

		go experiment(n, func() {
			runSimple(vm)
		})
	})

	mux.HandleFunc("GET /squirrel/{n}", func(w http.ResponseWriter, r *http.Request) {
		nStr := r.PathValue("n")

		n, err := strconv.Atoi(nStr)
		if err != nil {
			// Invalid integer → return 400 Bad Request
			http.Error(w, "n must be an integer", http.StatusBadRequest)
			return
		}

		go experiment(n, func() {
			runSquirrel(vm)
		})
	})

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	s := &http.Server{
		Handler: mux,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}
	go s.ListenAndServe()

	// Stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fsc.Stop()
}

func experiment(n int, f func()) {
	cfg := benchmark.Config{
		NumWorker:       n,
		Duration:        30 * time.Second,
		WarumupDuration: 5 * time.Second,
		MaxSleep:        0,
	}

	benchmark.Benchmark(cfg, f)
}

func runSquirrel(vm *viewregistry.Manager) {
	p := &simpleviews.SquirrelParams{}

	in, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	f, err := vm.NewView("squirrel", in)
	if err != nil {
		log.Fatal(err)
	}

	_ = f

	_, err = vm.InitiateView(f, context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func runSimple(vm *viewregistry.Manager) {
	p := &simpleviews.SimpleParams{}

	in, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}

	f, err := vm.NewView("simple", in)
	if err != nil {
		log.Fatal(err)
	}

	_ = f

	_, err = vm.InitiateView(f, context.Background())
	if err != nil {
		log.Fatal(err)
	}
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

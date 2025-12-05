/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/hyperledger-labs/fabric-smart-client/integration"
	"github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/api"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/cmd"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/cmd/network"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fabric"
	nwofabricx "github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fabricx"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fabricx/extensions/scv2"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fsc"
	nwofsc "github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fsc"
	"github.com/hyperledger-labs/fabric-smart-client/pkg/node"
	view "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/view/grpc/client/cmd"
	"github.com/onsi/gomega"
)

func main() {
	gomega.RegisterFailHandler(func(message string, callerSkip ...int) {
		panic(message)
	})

	m := cmd.NewMain("Simple network", "0.1")
	mainCmd := m.Cmd()
	network.StartCMDPostNew = func(infrastructure *integration.Infrastructure) error {
		infrastructure.RegisterPlatformFactory(nwofabricx.NewPlatformFactory())
		return nil
	}

	mainCmd.AddCommand(network.NewCmd(Topology(&simple.SDK{}, nwofsc.WebSocket)...))

	mainCmd.AddCommand(view.NewCmd())
	m.Execute()
}

func Topology(sdk node.SDK, commType fsc.P2PCommunicationType) []api.Topology {
	fabricTopology := nwofabricx.NewDefaultTopology()
	fabricTopology.AddOrganizationsByName("Org1")
	fabricTopology.AddNamespaceWithUnanimity("simple", "Org1")

	fscTopology := fsc.NewTopology()
	fscTopology.P2PCommunicationType = commType
	fscTopology.SetLogging("grpc=error:info", "")

	//fscTopology.AddNodeByName("simple").
	//	AddOptions(fabric.WithOrganization("Org1")).
	//	// simple responder
	//	RegisterResponder(&simpleviews.FoxView{}, &simpleviews.SquirrelView{})
	////RegisterResponder(&simpleviews.SimpleView{}, &simpleviews.RelayView{})

	fscTopology.AddNodeByName("perf").
		AddOptions(fabric.WithOrganization("Org1")).
		AddOptions(scv2.WithApproverRole()).
		SetExecutable("github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple/qnode/cmd").
		AddOptions(fsc.WithAlias("perfAlias"))

	fscTopology.AddNodeByName("simple").
		AddOptions(fabric.WithOrganization("Org1")).
		SetExecutable("github.com/hyperledger-labs/fabric-smart-client/integration/fabricx/simple/pnode/cmd")

	fscTopology.AddSDK(sdk)

	return []api.Topology{
		fabricTopology,
		fscTopology,
	}
}

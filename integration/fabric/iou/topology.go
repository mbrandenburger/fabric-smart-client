/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package iou

import (
	"github.com/hyperledger-labs/fabric-smart-client/integration/fabric/iou/views"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/api"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fabric"
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fsc"
)

func Topology() []api.Topology {
	// Define a Fabric topology with:
	// 1. Three organization: Org1, Org2, and Org3
	// 2. A namespace whose changes can be endorsed by Org1.
	fabricTopology := fabric.NewDefaultTopology()
	fabricTopology.AddOrganizationsByName("Org1", "Org2", "Org3")
	fabricTopology.SetNamespaceApproverOrgs("Org1")
	fabricTopology.AddNamespaceWithUnanimity("iou", "Org1")

	// Define an FSC topology with 3 FCS nodes.
	// One for the approver, one for the borrower, and one for the lender.
	fscTopology := fsc.NewTopology()

	// Add the approver FSC node.
	approver := fscTopology.AddNodeByName("approver")
	// This option equips the approver's FSC node with an identity belonging to Org1.
	// Therefore, the approver is an endorser of the Fabric namespace we defined above.
	approver.AddOptions(
		fabric.WithOrganization("Org1"),
	)
	approver.RegisterResponder(&views.ApproverView{}, &views.CreateIOUView{})
	approver.RegisterResponder(&views.ApproverView{}, &views.UpdateIOUView{})

	// Add the borrower's FSC node
	borrower := fscTopology.AddNodeByName("borrower")
	borrower.AddOptions(
		fabric.WithOrganization("Org2"),
	)
	borrower.RegisterViewFactory("create", &views.CreateIOUViewFactory{})
	borrower.RegisterViewFactory("update", &views.UpdateIOUViewFactory{})
	borrower.RegisterViewFactory("query", &views.QueryViewFactory{})

	// Add the lender's FSC node
	lender := fscTopology.AddNodeByName("lender")
	lender.AddOptions(
		fabric.WithOrganization("Org3"),
	)
	lender.RegisterResponder(&views.CreateIOUResponderView{}, &views.CreateIOUView{})
	lender.RegisterResponder(&views.UpdateIOUResponderView{}, &views.UpdateIOUView{})
	lender.RegisterViewFactory("query", &views.QueryViewFactory{})

	return []api.Topology{fabricTopology, fscTopology}
}

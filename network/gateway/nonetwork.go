//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package gateway

// TODO: spans, metrics

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

func newNoNetwork(b *Base) *NoNetwork {
	return &NoNetwork{Base: b}
}

// NoNetwork initial state
type NoNetwork struct {
	*Base

	// isDiscovery bool
	// skip        uint32
}

func (g *NoNetwork) Run(ctx context.Context) {

	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), g.NodeKeeper.GetOrigin().ID())

	if len(discoveryNodes) == 0 {
		g.zeroBootstrap(ctx)
		// create complete network
		g.Gatewayer.SwitchState(insolar.CompleteNetworkState)
		return
	}

	// run bootstrap
	isDiscovery := network.OriginIsDiscovery(cert)
	log.Info(isDiscovery)
	//getAuth

	// TODO: isDiscovery shaffle or sort

	for _, n := range discoveryNodes {
		h, _ := host.NewHostN(n.GetHost(), *n.GetNodeRef())

		res, err := g.BootstrapRequester.Authorize(ctx, h, cert)
		if err != nil {
			inslogger.FromContext(ctx).Errorf("Error authorizing to host %s: %s", h.String(), err.Error())
			continue
		}

		// TODO: check majority rule
		// if res.DiscoveryCount

		whichState := func() insolar.NetworkState {
			if network.OriginIsDiscovery(cert) {
				// TODO check ephemeral pulse res.PulseNumber
				return insolar.DiscoveryBootstrap
			}
			return insolar.JoinerBootstrap
		}

		switch res.Code {
		case packet.Success:
			g.permit = res.Permit
			g.Gatewayer.SwitchState(whichState())
			return
		case packet.WrongMandate:
			panic("wrong mandate " + res.Error)
		case packet.WrongVersion:
			panic("wrong mandate " + res.Error)
		}

	}

	// run again ?
	g.Gatewayer.SwitchState(insolar.NoNetworkState)

	// wait result from channel
	// save permits and change state

}

func (g *NoNetwork) GetState() insolar.NetworkState {
	return insolar.NoNetworkState
}

func (g *NoNetwork) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	return g.Base.OnPulse(ctx, pu)
}

func (g *NoNetwork) ShoudIgnorePulse(ctx context.Context, newPulse insolar.Pulse) bool {
	// if true { //!g.Base.NodeKeeper.IsBootstrapped() { always true here
	// 	g.Bootstrapper.SetLastPulse(newPulse.NextPulseNumber)
	// 	return true
	// }
	//
	// return g.isDiscovery && !g.NodeKeeper.GetConsensusInfo().IsJoiner() &&
	// 	newPulse.PulseNumber <= g.Bootstrapper.GetLastPulse()+insolar.PulseNumber(g.skip)
	// TODO: ??
	return true
}

// func (g *NoNetwork) connectToNewNetwork(ctx context.Context, address string) {
// 	g.NodeKeeper.GetClaimQueue().Push(&packets.ChangeNetworkClaim{Address: address})
// 	logger := inslogger.FromContext(ctx)
//
// 	// node, err := findNodeByAddress(address, g.CertificateManager.GetCertificate().GetDiscoveryNodes())
// 	// if err != nil {
// 	// 	logger.Warnf("Failed to find a discovery node: ", err)
// 	// }
//
// 	err := g.Bootstrapper.AuthenticateToDiscoveryNode(ctx, nil /*node*/)
// 	if err != nil {
// 		logger.Errorf("Failed to authenticate a node: " + err.Error())
// 	}
// }

// todo: remove this workaround
func (g *NoNetwork) zeroBootstrap(ctx context.Context) {
	inslogger.FromContext(ctx).Info("[ Bootstrap ] Zero bootstrap")
	g.NodeKeeper.SetInitialSnapshot([]insolar.NetworkNode{g.NodeKeeper.GetOrigin()})
}

func findNodeByAddress(address string, nodes []insolar.DiscoveryNode) (insolar.DiscoveryNode, error) {
	for _, node := range nodes {
		if node.GetHost() == address {
			return node, nil
		}
	}
	return nil, errors.New("Failed to find a discovery node with address: " + address)
}

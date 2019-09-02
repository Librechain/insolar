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

import (
	"context"
	"testing"

	"github.com/insolar/insolar/certificate"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/reply"

	"github.com/insolar/insolar/network"
	testnet "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/insolar"
)

func emtygateway(t *testing.T) network.Gateway {
	// todo use mockPulseManager(t)
	return newNoNetwork(&Base{})
}

func TestSwitch(t *testing.T) {
	t.Skip("fixme")
	ctx := context.Background()

	// nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper.MoveSyncToActiveMock.Set(func(ctx context.Context, number insolar.PulseNumber) {})
	gatewayer := testnet.NewGatewayerMock(t)
	// pm := mockPulseManager(t)

	ge := emtygateway(t)

	require.NotNil(t, ge)
	require.Equal(t, "NoNetworkState", ge.GetState().String())

	ge.Run(ctx, *insolar.EphemeralPulse)

	gatewayer.GatewayMock.Set(func() (g1 network.Gateway) {
		return ge
	})
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		ge = ge.NewGateway(ctx, state)
	})
	gilreleased := false

	ge.OnPulseFromPulsar(ctx, insolar.Pulse{}, nil)

	require.Equal(t, "CompleteNetworkState", ge.GetState().String())
	require.False(t, gilreleased)
	cref := gen.Reference()

	for _, state := range []insolar.NetworkState{insolar.NoNetworkState,
		insolar.JoinerBootstrap, insolar.CompleteNetworkState} {
		ge = ge.NewGateway(ctx, state)
		require.Equal(t, state, ge.GetState())
		ge.Run(ctx, *insolar.EphemeralPulse)
		au := ge.Auther()

		_, err := au.GetCert(ctx, &cref)
		require.Error(t, err)

		_, err = au.ValidateCert(ctx, &certificate.Certificate{})
		require.Error(t, err)

		ge.OnPulseFromPulsar(ctx, insolar.Pulse{}, nil)

	}

}

func TestDumbComplete_GetCert(t *testing.T) {
	t.Skip("fixme")
	ctx := context.Background()

	// nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper.MoveSyncToActiveMock.Set(func(ctx context.Context, number insolar.PulseNumber) {})

	gatewayer := testnet.NewGatewayerMock(t)

	CR := testutils.NewContractRequesterMock(t)
	CM := testutils.NewCertificateManagerMock(t)
	ge := emtygateway(t)
	// pm := mockPulseManager(t)

	// ge := newNoNetwork(gatewayer, pm,
	//	nodekeeper, CR,
	//	testutils.NewCryptographyServiceMock(t),
	//	testnet.NewHostNetworkMock(t),
	//	CM)

	require.NotNil(t, ge)
	require.Equal(t, "NoNetworkState", ge.GetState().String())

	ge.Run(ctx, *insolar.EphemeralPulse)

	gatewayer.GatewayMock.Set(func() (r network.Gateway) { return ge })
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		ge = ge.NewGateway(ctx, state)
	})
	gilreleased := false

	ge.OnPulseFromPulsar(ctx, insolar.Pulse{}, nil)

	require.Equal(t, "CompleteNetworkState", ge.GetState().String())
	require.False(t, gilreleased)

	cref := gen.Reference()

	CR.SendRequestMock.Set(func(ctx context.Context, ref *insolar.Reference, method string, argsIn []interface{}, p insolar.PulseNumber,
	) (r insolar.Reply, r2 *insolar.Reference, r1 error) {
		require.Equal(t, &cref, ref)
		require.Equal(t, "GetNodeInfo", method)
		repl, _ := insolar.Serialize(struct {
			PublicKey string
			Role      insolar.StaticRole
		}{"LALALA", insolar.StaticRoleVirtual})
		return &reply.CallMethod{
			Result: repl,
		}, nil, nil
	})

	CM.GetCertificateMock.Set(func() (r insolar.Certificate) { return &certificate.Certificate{} })
	cert, err := ge.Auther().GetCert(ctx, &cref)

	require.NoError(t, err)
	require.NotNil(t, cert)
	require.Equal(t, cert, &certificate.Certificate{})
}

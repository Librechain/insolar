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

package nodenetwork

import (
	"crypto"
	"github.com/insolar/insolar/network/storage"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNodeNetwork(t *testing.T) {
	cfg := configuration.Transport{Address: "invalid"}
	certMock := testutils.NewCertificateMock(t)
	certMock.GetRoleMock.Set(func() insolar.StaticRole { return insolar.StaticRoleUnknown })
	certMock.GetPublicKeyMock.Set(func() crypto.PublicKey { return nil })
	certMock.GetNodeRefMock.Set(func() *insolar.Reference { ref := gen.Reference(); return &ref })
	certMock.GetDiscoveryNodesMock.Set(func() []insolar.DiscoveryNode { return nil })
	_, err := NewNodeNetwork(cfg, certMock)
	assert.Error(t, err)
	cfg.Address = "127.0.0.1:3355"
	_, err = NewNodeNetwork(cfg, certMock)
	assert.NoError(t, err)
}

func newNodeKeeper(t *testing.T, service insolar.CryptographyService) network.NodeKeeper {
	cfg := configuration.Transport{Address: "127.0.0.1:3355"}
	certMock := testutils.NewCertificateMock(t)
	keyProcessor := platformpolicy.NewKeyProcessor()
	secret, err := keyProcessor.GeneratePrivateKey()
	require.NoError(t, err)
	pk := keyProcessor.ExtractPublicKey(secret)
	if service == nil {
		service = cryptography.NewKeyBoundCryptographyService(secret)
	}
	require.NoError(t, err)
	certMock.GetRoleMock.Set(func() insolar.StaticRole { return insolar.StaticRoleUnknown })
	certMock.GetPublicKeyMock.Set(func() crypto.PublicKey { return pk })
	certMock.GetNodeRefMock.Set(func() *insolar.Reference { ref := gen.Reference(); return &ref })
	certMock.GetDiscoveryNodesMock.Set(func() []insolar.DiscoveryNode { return nil })
	nw, err := NewNodeNetwork(cfg, certMock)
	require.NoError(t, err)
	nw.(*nodekeeper).SnapshotStorage = storage.NewMemoryStorage()
	return nw.(network.NodeKeeper)
}

func TestNewNodeKeeper(t *testing.T) {
	nk := newNodeKeeper(t, nil)
	origin := nk.GetOrigin()
	assert.NotNil(t, origin)
	nk.SetInitialSnapshot([]insolar.NetworkNode{origin})
	assert.NotNil(t, nk.GetAccessor(insolar.GenesisPulse.PulseNumber))
}

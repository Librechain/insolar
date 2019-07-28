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

// +build networktest

package tests

import (
	"fmt"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type bootstrapSuite struct {
	testSuite
}

func (s *bootstrapSuite) SetupTest() {
	s.fixtureMap[s.T().Name()] = newFixture(s.T())
	var err error
	s.fixture().pulsar, err = NewTestPulsar(reqTimeoutMs*10, pulseDelta*10)
	s.Require().NoError(err)

	inslogger.FromContext(s.fixture().ctx).Info("SetupTest -- ")

	for i := 0; i < s.bootstrapCount; i++ {
		s.fixture().bootstrapNodes = append(s.fixture().bootstrapNodes, s.newNetworkNode(fmt.Sprintf("bootstrap_%d", i)))
	}

	s.SetupNodesNetwork(s.fixture().bootstrapNodes)
	//for _, n := range s.fixture().bootstrapNodes {
	//	n.serviceNetwork.PhaseManager.(*phaseManagerWrapper).original = &fakeConsensus{n, s}
	//
	//}

	//pulseReceivers := make([]string, 0)
	//for _, node := range s.fixture().bootstrapNodes {
	//	pulseReceivers = append(pulseReceivers, node.host)
	//}
	//
	//log.Info("Start test pulsar")
	//err = s.fixture().pulsar.Start(s.fixture().ctx, pulseReceivers)
	//s.Require().NoError(err)
}

func (s *bootstrapSuite) TearDownTest() {
	inslogger.FromContext(s.fixture().ctx).Info("TearDownTest -- ")

	suiteLogger.Info("Stop bootstrap nodes")
	for _, n := range s.fixture().bootstrapNodes {
		err := n.componentManager.Stop(n.ctx)
		s.NoError(err)
	}
}

func (s *bootstrapSuite) waitForConsensus(consensusCount int) {
	for i := 0; i < consensusCount; i++ {
		for _, n := range s.fixture().bootstrapNodes {
			<-n.consensusResult
		}
	}
}

func newBootstraptSuite(bootstrapCount int) *bootstrapSuite {
	return &bootstrapSuite{
		testSuite: newTestSuite(bootstrapCount, 0),
	}
}

//func TestBootstrap(t *testing.T) {
//	s := newBootstraptSuite(2)
//	suite.Run(t, s)
//}

func (s *bootstrapSuite) TestExample() {
	inslogger.FromContext(s.fixture().ctx).Info("Log -- ")
	s.True(true)

	s.StartNodesNetwork(s.fixture().bootstrapNodes)

	s.waitForConsensus(1)

	//
	//
	// s.AssertActiveNodesCountDelta(0)
	//
	// s.waitForConsensus(1)
	//
	// s.AssertActiveNodesCountDelta(1)
	// s.AssertWorkingNodesCountDelta(0)
	//
	// s.waitForConsensus(2)
	//
	// s.AssertActiveNodesCountDelta(1)
	// s.AssertWorkingNodesCountDelta(1)
}

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

package syncrun

import (
	"context"
	"sync/atomic"
)

type ChainedCancel struct {
	state uint32 // atomic
	chain atomic.Value
}

func (p *ChainedCancel) Cancel() {
	for {
		lastState := atomic.LoadUint32(&p.state)
		switch {
		case lastState&0x01 != 0:
			return
		case !atomic.CompareAndSwapUint32(&p.state, lastState, lastState|0x01):
			continue
		case lastState == 0x04:
			p.runChain()
		}
		return
	}
}

func (p *ChainedCancel) runChain() {
	// here is a potential problem, because Go spec doesn't provide ANY ordering on atomic operations
	// but Go compiler does provide some guarantees, so lets hope for the best

	fn := (p.chain.Load()).(context.CancelFunc)
	if fn == nil {
		// this can only happen when atomic ordering is broken
		panic("unexpected atomic ordering")
	}
	fn()
}

func (p *ChainedCancel) IsCancelled() bool {
	return atomic.LoadUint32(&p.state)&0x01 != 0
}

func (p *ChainedCancel) SetChain(chain context.CancelFunc) {
	if chain == nil {
		panic("illegal value")
	}
	for {
		lastState := atomic.LoadUint32(&p.state)
		switch {
		case lastState&^0x01 != 0:
			return
		case !atomic.CompareAndSwapUint32(&p.state, lastState, lastState|0x02):
			continue
		}
		break
	}

	p.chain.Store(chain)

	for {
		lastState := atomic.LoadUint32(&p.state)
		switch {
		case lastState&^0x01 != 0x02:
			// this can only happen when atomic ordering is broken
			panic("unexpected atomic ordering")
		case !atomic.CompareAndSwapUint32(&p.state, lastState, (lastState&0x01)|0x04):
			continue
		case lastState&0x01 != 0:
			// if cancel was set then call the chained cancel here
			// otherwise, the cancelling process will be responsible to call the chained cancel
			p.runChain()
		}
		return
	}
}

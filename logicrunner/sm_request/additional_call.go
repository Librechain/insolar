//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package sm_request

import (
	"context"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/payload"
)

type StateMachineAdditionalCall struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.AdditionalCallFromPreviousExecutor
}

var declAdditionalCall smachine.StateMachineDeclaration = declarationAdditionalCall{}

type declarationAdditionalCall struct{}

func (declarationAdditionalCall) GetStepLogger(context.Context, smachine.StateMachine) smachine.StateMachineStepLoggerFunc {
	return nil
}

func (declarationAdditionalCall) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	_ = sm.(*StateMachineAdditionalCall)
}

func (declarationAdditionalCall) IsConsecutive(cur, next smachine.StateFunc) bool {
	return false
}

func (declarationAdditionalCall) GetShadowMigrateFor(smachine.StateMachine) smachine.ShadowMigrateFunc {
	return nil
}

func (declarationAdditionalCall) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineAdditionalCall)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachineAdditionalCall) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declAdditionalCall
}

func (s *StateMachineAdditionalCall) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Stop()
}
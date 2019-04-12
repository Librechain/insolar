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

package jet

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

type calculator struct {
	jetCoordinator insolar.JetCoordinator
	jetAccessor    Accessor
}

func NewCalculator(jetCoordinator insolar.JetCoordinator, jetAccessor Accessor) Calculator {
	return &calculator{jetCoordinator: jetCoordinator, jetAccessor: jetAccessor}
}

func (c *calculator) MineForPulse(ctx context.Context, pn insolar.PulseNumber) []insolar.JetID {
	var res []insolar.JetID

	jetIDs := c.jetAccessor.All(ctx, pn)
	me := c.jetCoordinator.Me()

	for _, jetID := range jetIDs {
		executor, err := c.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(jetID), pn)
		if err != nil {
			continue
		}
		if *executor == me {
			res = append(res, jetID)
		}
	}

	return res
}

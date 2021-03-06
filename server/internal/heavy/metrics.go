// Copyright 2020 Insolar Network Ltd.
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

package heavy

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statBadgerStartTime = stats.Float64(
		"badger_start_time",
		"Time of last badger starting",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statBadgerStartTime.Name(),
			Description: statBadgerStartTime.Description(),
			Measure:     statBadgerStartTime,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}

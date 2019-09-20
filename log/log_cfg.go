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

package log

import (
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

type parsedLogConfig struct {
	OutputType insolar.LogOutput
	LogFormat  insolar.LogFormat
	LogLevel   insolar.LogLevel

	SkipFrameBaselineAdjustment int
}

func parseLogConfig(cfg configuration.Log) (pc parsedLogConfig, err error) {

	pc.OutputType, err = insolar.ParseOutput(cfg.OutputType, insolar.DefaultLogOutput)
	if err != nil {
		return
	}

	pc.LogFormat, err = insolar.ParseFormat(cfg.Formatter, insolar.DefaultLogFormat)
	if err != nil {
		return
	}

	pc.LogLevel, err = insolar.ParseLevel(cfg.Level)
	if err != nil {
		return
	}

	return
}

//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package pulse

type Epoch uint32

const (
	InvalidPulseEpoch Epoch = iota

	// Epoch for pulses autogenerated by network during ephemeral network stage.
	EphemeralPulseEpoch

	// Epoch for missing pulses inside a range. Compatible with time pulses.
	ArticulationPulseEpoch

	MaxSpecialEpoch = iota - 1
	// NB! Special epoch is ONLY allowed to be [1..255], range of [256..65535] is FORBIDDEN
)

func (v Epoch) IsUnknown() bool {
	return v == InvalidPulseEpoch
}

func (v Epoch) IsEphemeral() bool {
	return v == EphemeralPulseEpoch
}

func (v Epoch) IsArticulation() bool {
	return v == ArticulationPulseEpoch
}

func (v Epoch) IsTimeEpoch() bool {
	isValid, isSpecial := v.Classify()
	return isValid && !isSpecial
}

func (v Epoch) Classify() (isValid, isSpecial bool) {
	switch {
	case v > MaxTimePulse:
		break
	case v >= MinTimePulse:
		return true, false
	case v <= MaxSpecialEpoch:
		return v != InvalidPulseEpoch, true
	}
	return false, false
}

func (v Epoch) IsValidEpoch() bool {
	isValid, _ := v.Classify()
	return isValid
}

func (v Epoch) IsCompatible(vv Epoch) bool {
	switch isValid, isSpecial := v.Classify(); {
	case !isValid:
		return false
	case !isSpecial:
		return vv == ArticulationPulseEpoch || vv.IsTimeEpoch()
	case v == vv:
		return true
	case v == ArticulationPulseEpoch:
		return vv.IsTimeEpoch()
	default:
		return false
	}
}

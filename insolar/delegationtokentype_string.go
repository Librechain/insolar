// Code generated by "stringer -type=DelegationTokenType"; DO NOT EDIT.

package insolar

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DTTypePendingExecution-1]
	_ = x[DTTypeGetObjectRedirect-2]
	_ = x[DTTypeGetCodeRedirect-3]
}

const _DelegationTokenType_name = "DTTypePendingExecutionDTTypeGetObjectRedirectDTTypeGetCodeRedirect"

var _DelegationTokenType_index = [...]uint8{0, 22, 45, 66}

func (i DelegationTokenType) String() string {
	i -= 1
	if i >= DelegationTokenType(len(_DelegationTokenType_index)-1) {
		return "DelegationTokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _DelegationTokenType_name[_DelegationTokenType_index[i]:_DelegationTokenType_index[i+1]]
}

// Code generated by "stringer -type=CommonSignal"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SUBTK_SUCCESS-0]
	_ = x[SUBTK_ERROR-1]
	_ = x[SUBTK_INIT-2]
	_ = x[SUBTK_DUPLICATED-3]
}

const _CommonSignal_name = "SUBTK_SUCCESSSUBTK_ERRORSUBTK_INITSUBTK_DUPLICATED"

var _CommonSignal_index = [...]uint8{0, 13, 24, 34, 50}

func (i CommonSignal) String() string {
	if i >= CommonSignal(len(_CommonSignal_index)-1) {
		return "CommonSignal(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CommonSignal_name[_CommonSignal_index[i]:_CommonSignal_index[i+1]]
}

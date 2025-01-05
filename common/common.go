package common

type CommonSignal uint

const (
	SUBTK_SUCCESS CommonSignal = iota
	SUBTK_ERROR
	SUBTK_INIT
	SUBTK_DUPLICATED
)
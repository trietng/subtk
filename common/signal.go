package common

type CommonSignal uint

const (
	SUBTK_SUCCESS CommonSignal = iota
	SUBTK_INIT
	SUBTK_DUPLICATED
	SUBTK_ERROR
)

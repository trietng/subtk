package signals

type ConfigSignal uint

const (
	Success ConfigSignal = iota
	Init
	Duplicated
	Error
)

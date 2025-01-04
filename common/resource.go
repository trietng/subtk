package common

type ResourceSignal uint

const (
	RESOURCE_OK ResourceSignal = iota
	RESOURCE_ERROR
	RESOURCE_NOT_FOUND
	RESOURCE_OUTDATED
)
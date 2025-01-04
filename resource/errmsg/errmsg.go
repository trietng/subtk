package errmsg

import "fmt"

var (
	// errors
	ErrInvalidResourceGroupName = fmt.Errorf("error: invalid resource group name")
	ErrInvalidResourceName      = fmt.Errorf("error: Invalid resource name")
)

func ErrGroupIntegrityCheckFailed(group string) error {
	return fmt.Errorf("error: integrity check failed for resource group '%s', run subtk repair -r to reset resources", group)
}

func ErrResourceIntegrityCheckFailed(group string, resource string) error {
	return fmt.Errorf("error: integrity check failed for resource '%s/%s', run subtk repair -r to reset resources", group, resource)
}

package errmsg

import (
	"fmt"
)

var (
	ErrNoSupportedFilesFound = fmt.Errorf("error: no supported files found")
)
package errmsg

import (
	"fmt"
)

var (
	ErrNoSupportedFilesFound = fmt.Errorf("error: no supported files found")
	ErrUnableToDetectEpisode = fmt.Errorf("error: unable to detect episode")
)
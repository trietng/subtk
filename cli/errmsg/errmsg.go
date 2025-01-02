package errmsg

import (
	"fmt"
)

var (
	ErrInvalidModule         = fmt.Errorf("error: invalid module")
	ErrInvalidApiKeyFormat   = fmt.Errorf("error: invalid api key format, should be in the form of <provider>:<api_key>")
	ErrFeatureNotImplemented = fmt.Errorf("error: feature not implemented")
)
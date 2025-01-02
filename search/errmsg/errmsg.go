package errmsg

import (
	"fmt"
	"trietng/subtk/search/datasources"
)

var (
	ErrInsufficientNumberOfDataSources = fmt.Errorf("error: insufficient number of suppliers to fetch data (requires at least 1)")
)

func WarnSearchRequestFailed(datasource datasources.DataSource) string {
	return fmt.Sprintf("warning: search request to %s failed", datasource.Endpoint())
}
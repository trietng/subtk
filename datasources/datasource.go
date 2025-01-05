package datasources

import "trietng/subtk/search/result"

type DataSource interface {
	Name() string
	Endpoint() string
	Search() ([]result.SubtitleSearchResult, error)
	VerifyDownloadUrl(url string) bool
}
package query

import "trietng/subtk/common"

type QueryType uint

const (
	QUERY_RELAXED      QueryType = iota // simple query, i.e. similar to a google search
	QUERY_RELEASE_NAME                  // query is a release name
)

type QueryMetadata struct {
	QueryType       QueryType        // type of query
	MediaType       common.MediaType // type of media
	HearingImpaired bool        // whether the subtitles are hearing impaired
}

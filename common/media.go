package common

type MediaType uint

const (
	NONE MediaType = iota
	MOVIE
	TV
)

type MediaInfo struct {
	Title   string
	Type    MediaType
	Season  int
	Episode int
}
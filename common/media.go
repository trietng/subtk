package common

type MediaType uint

const (
	MOVIE MediaType = iota
	TV
	OTHER
)

type MediaInfo struct {
	Title  string
	Type   MediaType
	Season int
}

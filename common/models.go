package common

type MediaType uint

const (
	Movie MediaType = iota
	TV
)

type MediaInfo struct {
	Title  string
	Type   MediaType
	Season int
	Year   int
}
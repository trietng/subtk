package common

type MediaType uint

const (
	Movie MediaType = iota
	TV
	Other
)

type MediaInfo struct {
	Title  string
	Type   MediaType
	Season int
}
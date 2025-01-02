package merger

import (
	"trietng/subtk/search/merger/mergestrategy"
	"trietng/subtk/search/result"
)

type Merger struct {
	Strategy string
}

func (m *Merger) Merge(data [][]result.SubtitleSearchResult) []result.SubtitleSearchResult {
	switch m.Strategy {
	case mergestrategy.First:
		if len(data) < 1 {
			return []result.SubtitleSearchResult{}
		}
		return data[0]
	default:
		return []result.SubtitleSearchResult{}
	}
}
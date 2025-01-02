package search

import (
	"fmt"
	"strings"
	"sync"
	"trietng/subtk/config"
	"trietng/subtk/search/datasources"
	"trietng/subtk/search/errmsg"
	"trietng/subtk/search/merger"
	"trietng/subtk/search/result"
)

type SubtitleSearchEngine struct {
	query         string
	dataSources   []datasources.DataSource
	mergeStrategy string
}

func NewSubtitleSearchEngine(query string, mergeStrategy string) *SubtitleSearchEngine {
	// load configuration on startup
	config.LoadConfig()
	// create data sources
	// 1. subdl
	subdl := datasources.NewSubdlDataSource(query)
	// register data sources
	dataSources := []datasources.DataSource{
		subdl,
	}
	// enforce lowercase on merge strategy
	mergeStrategy = strings.ToLower(mergeStrategy)
	// return the search engine
	return &SubtitleSearchEngine{
		query: query,
		dataSources: dataSources,
		mergeStrategy: mergeStrategy,
	}
}

func (se *SubtitleSearchEngine) Search(query string) ([]result.SubtitleSearchResult, error) {
	n := len(se.dataSources)
	if n < 1 {
		return nil, errmsg.ErrInsufficientNumberOfDataSources
	}
	// create a channel to receive the results
    channel := make(chan []result.SubtitleSearchResult, n)
    // wait for all goroutines to finish
    var wg sync.WaitGroup
    for _, dataSource := range se.dataSources {
        wg.Add(1)
        // goroutine to fetch data from supplier
        go func(dataSource datasources.DataSource) {
            defer wg.Done()
            result, err := dataSource.Search()
            if err == nil {
                channel <- result
            } else {
                fmt.Println(err)
                fmt.Println(errmsg.WarnSearchRequestFailed(dataSource))
                channel <- nil
            }
        }(dataSource)
	}
	// finalize the results
	var data [][]result.SubtitleSearchResult
	for range se.dataSources {
		item := <-channel
		if item != nil {
			data = append(data, item)
		}
	}
	// retrieve merge strategy
	resultMerger := merger.Merger {
		Strategy: se.mergeStrategy,
	}
	// merge the results
	merged := resultMerger.Merge(data)
	return merged, nil
}
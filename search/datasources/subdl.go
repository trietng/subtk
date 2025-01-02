package datasources

import (
	"encoding/json"
	"io"
	"regexp"
	"slices"
	"strconv"
	"sync"

	//"fmt"
	"net/http"
	"net/url"
	"trietng/subtk/config"

	//"trietng/subtk/search"
	"trietng/subtk/search/result"
)

var (
	// regex to extract the download count from the page, e.g. "downloads":123
	downloadCountRegex = regexp.MustCompile(`"downloads":(\d+)`)
)

type SubdlDataSource struct {
	apiKey      string
	queryValues url.Values
}

func (ds *SubdlDataSource) Name() string {
	return "subdl"
}

func (ds *SubdlDataSource) Endpoint() string {
	return "https://api.subdl.com/api/v1/subtitles"
}

func (ds *SubdlDataSource) buildInfoUrl(rawUrl string) string {
	return "https://subdl.com" + rawUrl
}

func (ds *SubdlDataSource) buildDownloadUrl(rawUrl string) string {
	return "https://dl.subdl.com" + rawUrl
}

func (ds *SubdlDataSource) Search() ([]result.SubtitleSearchResult, error) {
	apiKeyQueryValue := url.Values{
		"api_key": {ds.apiKey},
	}
	resp, err := http.Get(ds.Endpoint() + "?" + apiKeyQueryValue.Encode() + "&" + ds.queryValues.Encode())
	if err != nil {
		return nil, err
	}
	var object map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&object)
    if err != nil {
        return nil, err
    }
	defer resp.Body.Close()
	if subtitles, ok := object["subtitles"].([]interface{}); ok {
		searchResults := make([]result.SubtitleSearchResult, len(subtitles))
		for index, unmappedSubtitle := range subtitles {
			if subtitle, ok := unmappedSubtitle.(map[string]interface{}); ok {
				searchResults[index].Provider = ds.Name()
				if releaseName, ok := subtitle["release_name"].(string); ok {
					searchResults[index].ReleaseName = releaseName
				}
				if author, ok := subtitle["author"].(string); ok {
					searchResults[index].Author = author
				}
				if url, ok := subtitle["url"].(string); ok {
					searchResults[index].DownloadUrl = ds.buildDownloadUrl(url)
				}
				if subtitlePage, ok := subtitle["subtitlePage"].(string); ok {
					searchResults[index].InfoUrl = ds.buildInfoUrl(subtitlePage)
				}
			}
		}
		// score the results
		// criteria: download counts
		//scores := make([]int, len(searchResults))
		channel := make(chan int, len(searchResults))
		var wg sync.WaitGroup
		for _, searchResult := range searchResults {
			wg.Add(1)
			defer wg.Done()
			go func(searchResult *result.SubtitleSearchResult) {
				// download the page
				resp, err := http.Get(searchResult.InfoUrl)
				if err != nil {
					channel <- -1
				} else {
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						channel <- -1
					} else {
						// extract the download count
						matches := downloadCountRegex.FindSubmatch(body)
						if len(matches) < 2 {
							channel <- -1
						} else {
							downloadCount := string(matches[1])
							score, err := strconv.Atoi(downloadCount)
							if err != nil {
								channel <- -1
							} else {
								channel <- score
							}
						}
						channel <- -1
					}
				}
			}(&searchResult)
		}
		for index := range searchResults {
			searchResults[index].Score = <-channel
		}
		// sort the results
		slices.SortStableFunc(searchResults, func(a, b result.SubtitleSearchResult) int {
			return b.Score - a.Score
		})
		return searchResults, nil
	}
	return []result.SubtitleSearchResult{}, nil
}

func NewSubdlDataSource(query string) *SubdlDataSource {
	ds := &SubdlDataSource{
		queryValues: url.Values{},
	}
	if key, ok := config.GetApiKey(ds.Name()); ok {
		ds.apiKey = key
	}
	ds.queryValues.Add("file_name", query)
	ds.queryValues.Add("languages", "EN")
	return ds
}
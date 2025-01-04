package datasources

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"trietng/subtk/config"
	"trietng/subtk/resource/languages"
	"trietng/subtk/resource/languages/fallback"
	"trietng/subtk/search/query"
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
					}
				}
			}(&searchResult)
		}
		for index := range searchResults {
			searchResults[index].Score = <-channel
		}
		// sort the results
		// stable sort due to concurrency
		slices.SortStableFunc(searchResults, func(a, b result.SubtitleSearchResult) int {
			return b.Score - a.Score
		})
		return searchResults, nil
	}
	return []result.SubtitleSearchResult{}, nil
}

func NewSubdlDataSource(q string, metadata query.QueryMetadata) *SubdlDataSource {
	ds := &SubdlDataSource{
		queryValues: url.Values{},
	}
	if key, ok := config.GetApiKey(ds.Name()); ok {
		ds.apiKey = key
	}
	supportedLanguages := (languages.SubdlLanguagesRepository{}).GetSupportedLanguages()
	userDefaultLanguages := config.GetDefaultLanguage()
	defaultLanguage := strings.ToUpper(fallback.DefaultLanguage)
	if _, ok := supportedLanguages[userDefaultLanguages]; ok {
		defaultLanguage = strings.ToUpper(userDefaultLanguages)
	}
	ds.queryValues.Add("languages", defaultLanguage)
	if (metadata.QueryType == query.QUERY_RELEASE_NAME) {
		ds.queryValues.Add("file_name", q)
	} else {
		ds.queryValues.Add("film_name", q)
	}
	return ds
}
package match

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"trietng/subtk/common"
	"trietng/subtk/match/errmsg"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

type SubtitleMatcher struct{}

var (
	SupportedFileTypes = map[string]struct{}{
		".avi":  {},
		".m4v":  {},
		".mkv":  {},
		".mp4":  {},
		".mpg":  {},
		".mpeg": {},
	}
)

func (m *SubtitleMatcher) Report() (*common.MediaInfo, error) {
	// get all files in the current directory that have a supported file type
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var report *common.MediaInfo
	var count int
	for _, file := range files {
		// skip directories
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if _, ok := SupportedFileTypes[ext]; ok {
			if report == nil {
				report = &common.MediaInfo{}
				report.Title = strings.TrimSuffix(file.Name(), ext)
				matches := common.REGEXP_SEASON.FindStringSubmatch(report.Title)
				if len(matches) > 1 {
					report.Season, err = strconv.Atoi(matches[1])
					if err != nil {
						return nil, err
					}
					report.Type = common.TV
				} else {
					report.Type = common.MOVIE
				}
			}
			count++
		}
	}
	// if there is more than one supported file, check for directory name
	if count > 1 {
		// get current folder name
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = filepath.Base(dir)
		similarity := strutil.Similarity(dir, report.Title, metrics.NewLevenshtein())
		if similarity > 0.4 {
			report.Title = dir
		} // else the dir name is bogus and we ignore it
	}
	if report == nil {
		return nil, errmsg.ErrNoSupportedFilesFound
	}
	return report, nil
}

func (m *SubtitleMatcher) MatchFiles(dir string) ([]common.MediaInfo, error) {
	// get all files in the current directory that have a supported file type
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var mis []common.MediaInfo
	for _, file := range files {
		// skip directories
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if _, ok := SupportedFileTypes[ext]; ok {
			mi := common.MediaInfo{}
			mi.Title = strings.TrimSuffix(file.Name(), ext)
			mi.Season, mi.Episode, err = GetSeasonAndEpisode(mi.Title)
			if err == nil {
				mis = append(mis, mi)
			}
		}
	}
	return mis, nil
}

func Summary(mediaInfo *common.MediaInfo) string {
	return fmt.Sprintf("file name: %s\ntype: %s\nseason: %d", mediaInfo.Title, mediaInfo.Type, mediaInfo.Season)
}

func GetSeasonAndEpisode(title string) (int, int, error) {
	matches := common.REGEXP_SEASON_EPISODE.FindStringSubmatch(title)
	if len(matches) < 3 {
		return 0, 0, errmsg.ErrUnableToDetectEpisode
	}
	season, _ := strconv.Atoi(matches[common.REGEXP_INDEX_SEASON])
	episode, _ := strconv.Atoi(matches[common.REGEXP_INDEX_EPISODE])
	return season, episode, nil
}
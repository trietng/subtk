package match

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"trietng/subtk/common"
	"trietng/subtk/match/errmsg"
)

type Matcher struct{}

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

var (
	// regex to extract season from a file name
	seasonRegexp = regexp.MustCompile(`(?i)(?:season\s*|s)(\d+)(?:e\d+|x\d+|)`)
)

func (a *Matcher) Match() (*common.MediaInfo, error) {
	// get all files in the current directory that have a supported file type
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var report common.MediaInfo
	for _, file := range files {
		// skip directories
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if _, ok := SupportedFileTypes[ext]; ok {
			report.Title = strings.TrimSuffix(file.Name(), ext)
			matches := seasonRegexp.FindStringSubmatch(report.Title)
			if len(matches) > 1 {
				report.Season, err = strconv.Atoi(matches[1])
				if err != nil {
					return nil, err
				}
				report.Type = common.TV
			}
			break
		}
	}
	// if filename is empty, return an error
	if report.Title == "" {
		return nil, errmsg.ErrNoSupportedFilesFound
	}
	return &report, nil
}

func Summary(mediaInfo *common.MediaInfo) string {
	return fmt.Sprintf("file name: %s\ntype: %s\nseason: %d", mediaInfo.Title, mediaInfo.Type, mediaInfo.Season)
}
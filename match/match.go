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
	// regex to extract title, season and year from a file name
	// e.g. "The.Fall.of.the.House.of.Usher.(2023).S01E01.A.Midnight.Dreary.1080p.NF.WEB-DL.10bit.DDP5.1.Atmos.x265-YELLO"
	// or "The Fall of the House of Usher 2023 1x1"
	r = regexp.MustCompile(`(?i)(?:season\s*|s)(\d+)(?:e\d+|x\d+|)`)
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
			matches := r.FindStringSubmatch(report.Title)
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
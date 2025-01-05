package common

import "regexp"

var (
	// regex to extract season from a file name
	REGEXP_SEASON = regexp.MustCompile(`(?i)(?:season\s*|s)(\d+)[-\.\s]*(?:e\d+|x\d+|)`)
	// regex to extract season and episode from a file name
	REGEXP_SEASON_EPISODE = regexp.MustCompile(`(?i)(?:season\s*|s)(?P<season>\d+)[-\.\s]*(?:episode\s*|e|x)(?P<episode>\d+)`)
	REGEXP_INDEX_SEASON = REGEXP_SEASON_EPISODE.SubexpIndex("season")
	REGEXP_INDEX_EPISODE = REGEXP_SEASON_EPISODE.SubexpIndex("episode")
)
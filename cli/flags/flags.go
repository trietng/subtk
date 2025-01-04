package flags

import (
	"flag"
	"os"
	"trietng/subtk/cli/module"
)

// declare flags
var (
	// flags for each module
	ConfigFlags struct {
		ApiKeyList         *bool
		ApiKeySet          *string
		ApiKeyUnset        *string
		DefaultLanguageSet *string
	}
	RepairFlags struct {
		Resource *bool
	}
	SearchFlags struct {
		Query           *string
		MergeStrategy   *string
		HearingImpaired *bool
	}
	DownloadFlags struct {
		Query		   *string
		Url            *string
		Destination    *string
		ExtractArchive *bool
	}
)

// sets flags for the specified module
func SetModuleFlags(mod string) {
	switch mod {
	case module.Config:
		ConfigFlags.ApiKeyList = flag.Bool("al", false, "list all api keys")
		ConfigFlags.ApiKeySet = flag.String("as", "", "api key to set <provider>:<api_key>")
		ConfigFlags.ApiKeyUnset = flag.String("au", "", "api key to unset")
		ConfigFlags.DefaultLanguageSet = flag.String("dls", "", "default language to set, must be a valid ISO 639-1 language code")
	case module.Repair:
		RepairFlags.Resource = flag.Bool("r", false, "reset all resources")
	case module.Search:
		SearchFlags.Query = flag.String("q", "", "query to search; leave empty for auto search")
		SearchFlags.MergeStrategy = flag.String("ms", "first", "merge strategy to use when merging search results")
		SearchFlags.HearingImpaired = flag.Bool("hi", false, "whether to include hearing impaired subtitles")
	case module.Download:
		DownloadFlags.Query = flag.String("q", "", "query to search and download; leave empty for auto search")
		dir, err := os.Getwd()
		if err != nil {
			dir = ""
		}
		DownloadFlags.Destination = flag.String("d", dir, "destination path to download the subtitle to")
		DownloadFlags.ExtractArchive = flag.Bool("ea", true, "whether to extract the downloaded archive")
		DownloadFlags.Url = flag.String("u", "", "direct url of the subtitle to download")
	}
	flag.Parse()
}
package flags

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
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
		Url            *string
		Destination    *string
		ExtractArchive *bool
	}
)

// sets flags for the specified module
func SetModuleFlags(mod string) {
	switch mod {
	case module.Console:
		flag.Usage = func() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			defer w.Flush()
			fmt.Fprintln(w, "Usage: subtk <module> [flags]")
			fmt.Fprintln(w, "MODULE\tDESCRIPTION")
			for _, m := range module.Modules {
				fmt.Fprintf(w, "%s\t%s\n", m, module.ModuleDescriptions[m])
			}
		}
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
		SearchFlags.HearingImpaired = flag.Bool("hi", false, "whether to enforce hearing impaired subtitles")
	case module.Download:
		SearchFlags.Query = flag.String("q", "", "query to search and download; leave empty for auto search")
		SearchFlags.MergeStrategy = flag.String("ms", "first", "merge strategy to use when merging search results")
		SearchFlags.HearingImpaired = flag.Bool("hi", false, "whether to enforce hearing impaired subtitles")
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
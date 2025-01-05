package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"trietng/subtk/cli/errmsg"
	"trietng/subtk/cli/flags"
	"trietng/subtk/cli/module"
	"trietng/subtk/config"
	"trietng/subtk/download"
	"trietng/subtk/match"
	"trietng/subtk/repair"
	"trietng/subtk/search"
	"trietng/subtk/search/query"
)

func Run(mod string) {
	// make module lower case
	mod = strings.ToLower(mod)
	switch mod {
	case module.Config:
		if *flags.ConfigFlags.ApiKeyList {
			serialized, _ := json.MarshalIndent(config.GetApiKeys(), "", "  ")
			fmt.Println(string(serialized))
		} else if *flags.ConfigFlags.ApiKeySet != "" {
			// pre-process the input
			*flags.ConfigFlags.ApiKeySet = strings.TrimSpace(*flags.ConfigFlags.ApiKeySet)
			// split the input
			parts := strings.Split(*flags.ConfigFlags.ApiKeySet, ":")
			if len(parts) != 2 {
				fmt.Println(errmsg.ErrInvalidApiKeyFormat)
				return
			}
			// set the api key
			config.SetApiKey(parts[0], parts[1])
			config.SaveConfig()
		} else if *flags.ConfigFlags.ApiKeyUnset != "" {
			// pre-process the input
			*flags.ConfigFlags.ApiKeySet = strings.TrimSpace(*flags.ConfigFlags.ApiKeySet)
			// unset the api key
			config.UnsetApiKey(*flags.ConfigFlags.ApiKeyUnset)
			config.SaveConfig()
		} else {
			goto ManualHelp
		}
	case module.Repair:
		if *flags.RepairFlags.Resource {
			// reset all resources
			repair.ResetResources()
		} else {
			goto ManualHelp
		}
	case module.Match:
		analyzer := match.SubtitleMatcher{}
		report, err := analyzer.Report()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(match.Summary(report))
		}
	case module.Search:
		// search for subtitles
		if *flags.SearchFlags.Query != "" {
			searchEngine := search.NewSubtitleSearchEngine(
				*flags.SearchFlags.Query,
				*flags.SearchFlags.MergeStrategy,
				query.QueryMetadata{
					HearingImpaired: *flags.SearchFlags.HearingImpaired,
				},
			)
			results, err := searchEngine.Search(*flags.SearchFlags.Query)
			if err != nil {
				fmt.Println(err)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			if len(results) == 0 {
				fmt.Println("no results found")
			} else {
				fmt.Fprintln(w, "Provider\tRelease Name\tDownload Url\tScore")
				for _, result := range results {
					fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", result.Provider, result.ReleaseName, result.DownloadUrl, result.Score)
				}
				w.Flush()
			}
		} else {
			goto ManualHelp
		}
	case module.Download:
		if *flags.DownloadFlags.Url != "" {
			// download subtitles
			downloader := download.SubtitleDownloader{
				Dir:            *flags.DownloadFlags.Destination,
				ExtractArchive: *flags.DownloadFlags.ExtractArchive,
				Url:            *flags.DownloadFlags.Url,
			}
			err := downloader.Download()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// search for subtitles
			metadata := query.QueryMetadata{
				HearingImpaired: *flags.SearchFlags.HearingImpaired,
			}
			searchQuery := *flags.SearchFlags.Query
			// if no query is provided, engage in automatic mode
			// automatic mode is only available for the download module
			if *flags.SearchFlags.Query == "" {
				metadata.QueryType = query.QUERY_RELEASE_NAME
				matcher := match.SubtitleMatcher{}
				report, err := matcher.Report()
				if err != nil {
					fmt.Println(err)
					return
				}
				searchQuery = report.Title
			}
			searchEngine := search.NewSubtitleSearchEngine(
				searchQuery,
				*flags.SearchFlags.MergeStrategy,
				metadata,
			)
			results, err := searchEngine.Search(*flags.SearchFlags.Query)
			if err != nil {
				fmt.Println(err)
			} else {
				// download the first result
				if len(results) > 0 {
					downloader := download.SubtitleDownloader{
						Dir:            *flags.DownloadFlags.Destination,
						ExtractArchive: *flags.DownloadFlags.ExtractArchive,
						Url:            results[0].DownloadUrl,
					}
					err := downloader.Download()
					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("no results found")
				}
			}
		}
	default:
		fmt.Println(errmsg.ErrInvalidModule)
	}
	return
ManualHelp:
	fmt.Printf("Usage of %s:\n", mod)
	flag.PrintDefaults()
}

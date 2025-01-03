package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"trietng/subtk/cli/errmsg"
	"trietng/subtk/cli/flags"
	"trietng/subtk/cli/module"
	"trietng/subtk/config"
	"trietng/subtk/match"
	"trietng/subtk/search"
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
		}
	case module.Match:
		analyzer := match.Matcher{}
		report, err := analyzer.Match()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(match.Summary(report))
		}
	case module.Search:
		// search for subtitles
		searchEngine := search.NewSubtitleSearchEngine(
			*flags.SearchFlags.Query,
			*flags.SearchFlags.MergeStrategy,
		)
		results, err := searchEngine.Search(*flags.SearchFlags.Query)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(results)
	case module.Download:
		fmt.Println(errmsg.ErrFeatureNotImplemented)
	default:
		fmt.Println(errmsg.ErrInvalidModule)
	}
}

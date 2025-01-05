package module

var (
	Modules = []string{
		"console",
		"config",
		"repair",
		"match",
		"download",
		"search",
	}
	Console  = Modules[0]
	Config   = Modules[1]
	Repair   = Modules[2]
	Match    = Modules[3]
	Download = Modules[4]
	Search   = Modules[5]
	// ModuleDescriptions is a map of module names to their descriptions
	ModuleDescriptions = map[string]string{
		Console:  "default module",
		Config:   "configure subtk",
		Repair:   "repair subtk",
		Match:    "match subtitles and generate report",
		Download: "download subtitles",
		Search:   "search for subtitles",
	}
)
package languages

import (
	"encoding/json"
	"net/http"
	"strings"
	"trietng/subtk/common"
	"trietng/subtk/common/utils/setops"
	"trietng/subtk/resource"
	"trietng/subtk/resource/languages/fallback"
	"trietng/subtk/resource/languages/iso639"
)

type SubdlLanguagesRepository struct{}

func (s SubdlLanguagesRepository) GetSupportedLanguages() map[string]struct{} {
	// check if resouce file for subdl supported languages exists
	// if not, create it
	if supportedLanguages, sig := resource.GetResource[map[string]struct{}]("languages", "subdl"); sig == common.RESOURCE_OK {
		return *supportedLanguages
	} else {
		// download supported languages from subdl
		resp, err := http.Get("https://subdl.com/api-files/language_list.json")
		if err != nil {
			return fallback.SupportedLanguages
		}
		defer resp.Body.Close()
		var temp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&temp)
		if err != nil {
			return fallback.SupportedLanguages
		}
		remoteLanguages := setops.Mtos(temp)
		// make all keys lower case
		for k := range remoteLanguages {
			delete(remoteLanguages, k)
			remoteLanguages[strings.ToLower(k)] = struct{}{}
		}
		// intersect with iso639-1
		intersectedLanguages := setops.Intersect(remoteLanguages, setops.Mtos(iso639.Set1))
		// save to resource file
		resource.SetResource("languages", "subdl", intersectedLanguages, resource.UPDATE_DEFAULT)
		return intersectedLanguages
	}
}

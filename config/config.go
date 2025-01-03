package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"trietng/subtk/common"
)

var (
	// config map
	config map[string]interface{}
	lock   sync.Mutex
	// error messages
	errorUnableToLoadConfigFile = fmt.Errorf("error: Unable to load config file")
	errorUnableToSaveConfigFile = fmt.Errorf("error: Unable to save config file")
)

const (
	apikeysKey         = "apikeys"
	defaultLanguageKey = "default_language"
)

func init() {
	LoadConfig()
}

func LoadConfig() (common.CommonSignal, error) {
	lock.Lock()
	defer lock.Unlock()
	var configSignal common.CommonSignal
	if config != nil {
		configSignal = common.SUBTK_DUPLICATED
		return configSignal, nil
	}
	home, _ := os.UserHomeDir()
	// create .subtk directory if it does not exist
	os.MkdirAll(fmt.Sprintf("%s/.subtk", home), 0755)
	file, err := os.Open(fmt.Sprintf("%s/.subtk/config.json", home))
	// if the file does not exist, create it
	if err != nil {
		os.WriteFile(fmt.Sprintf("%s/.subtk/config.json", home), []byte("{}"), 0644)
		file, _ = os.Open(fmt.Sprintf("%s/.subtk/config.json", home))
		configSignal = common.SUBTK_INIT
	}
	defer file.Close()
	raw, err := io.ReadAll(file)
	if err != nil {
		goto ConfigError
	}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		goto ConfigError
	}
	return configSignal, nil
ConfigError:
	configSignal = common.SUBTK_ERROR
	return configSignal, errorUnableToLoadConfigFile
}

func SaveConfig() error {
	lock.Lock()
	defer lock.Unlock()
	home, _ := os.UserHomeDir()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errorUnableToSaveConfigFile
	}
	os.WriteFile(fmt.Sprintf("%s/.subtk/config.json", home), data, 0644)
	return nil
}

func GetApiKeys() map[string]interface{} {
	lock.Lock()
	defer lock.Unlock()
	if apikeys, ok := config[apikeysKey].(map[string]interface{}); ok {
		return apikeys
	}
	return map[string]interface{}{}
}

func GetApiKey(name string) (string, bool) {
	lock.Lock()
	defer lock.Unlock()
	var (
		apikeys map[string]interface{}
		ok      bool
		value   string
	)
	if apikeys, ok = config[apikeysKey].(map[string]interface{}); ok {
		value, ok = apikeys[name].(string)
	}
	return value, ok
}

func SetApiKey(name string, value string) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := config[apikeysKey]; !ok {
		config[apikeysKey] = make(map[string]interface{})
	}
	apikeys := config[apikeysKey].(map[string]interface{})
	apikeys[name] = value
}

func UnsetApiKey(name string) {
	lock.Lock()
	defer lock.Unlock()
	if apikeys, ok := config[apikeysKey].(map[string]interface{}); ok {
		delete(apikeys, name)
	}
}

func getConfig(key string) (interface{}, bool) {
	lock.Lock()
	defer lock.Unlock()
	value, ok := config[key]
	return value, ok
}

func setConfig(key string, value interface{}) {
	lock.Lock()
	defer lock.Unlock()
	config[key] = value
}

func SetDefaultLanguage(lang string) {
	setConfig(defaultLanguageKey, lang)
}

func GetDefaultLanguage(lang string) (string, bool) {
	if value, ok := getConfig(defaultLanguageKey); ok {
		return value.(string), ok
	}
	return "", false
}

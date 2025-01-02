package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"trietng/subtk/config/signals"
)

var (
	// config map
	config map[string]interface{}
	lock   sync.Mutex
	// warning messages
	warningConfigAlreadyLoaded  = "warning: Config already loaded"
	// error messages
	errorUnableToLoadConfigFile = fmt.Errorf("error: Unable to load config file")
	errorUnableToSaveConfigFile = fmt.Errorf("error: Unable to save config file")
)

func LoadConfig() (signals.ConfigSignal, error) {
	lock.Lock()
	defer lock.Unlock()
	var configSignal signals.ConfigSignal
	if config != nil {
		fmt.Println(warningConfigAlreadyLoaded)
		configSignal = signals.Duplicated
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
		configSignal = signals.Init
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
	configSignal = signals.Error
	return configSignal, errorUnableToLoadConfigFile
}

func SaveConfig() error {
	lock.Lock()
	defer lock.Unlock()
	home, _ := os.UserHomeDir()
	data, err := json.Marshal(config)
	if err != nil {
		return errorUnableToSaveConfigFile
	}
	os.WriteFile(fmt.Sprintf("%s/.subtk/config.json", home), data, 0644)
	return nil
}

func GetApiKey(name string) (string, bool) {
	lock.Lock()
	defer lock.Unlock()
	var (
		apikeys map[string]interface{}
		ok      bool
		value   string
	)
	if apikeys, ok = config["apikeys"].(map[string]interface{}); ok {
		value, ok = apikeys[name].(string)
	}
	return value, ok
}

func SetApiKey(name string, value string) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := config["apikeys"]; !ok {
		config["apikeys"] = make(map[string]interface{})
	}
	apikeys := config["apikeys"].(map[string]interface{})
	apikeys[name] = value
}

func UnsetApiKey(name string) {
	lock.Lock()
	defer lock.Unlock()
	if apikeys, ok := config["apikeys"].(map[string]interface{}); ok {
		delete(apikeys, name)
	}
}

func GetApiKeys() map[string]interface{} {
	lock.Lock()
	defer lock.Unlock()
	if apikeys, ok := config["apikeys"].(map[string]interface{}); ok {
		return apikeys
	}
	return map[string]interface{}{}
}
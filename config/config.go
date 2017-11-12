//jasonxu
package config

import (
	"errors"
	"os"
	"encoding/json"
)

func GetConfigFromFile(path string, config interface{}) error {
	if path == "" {
		return errors.New("please check file path")
	}

	fileInfo, err := os.Open(path)

	if err != nil {
		return err
	}

	defer fileInfo.Close()

	jsonParser := json.NewDecoder(fileInfo)
	return jsonParser.Decode(&config)
}

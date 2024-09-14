package jsonfile

import (
	"encoding/json"
	"os"
)

func ReadFileToMap(Filename string, Map *map[string]string) error {
	file, err := os.OpenFile(Filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(Map)
}

func WriteMapToFile(Filename string, Map *map[string]string) error {
	file, err := os.OpenFile(Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(Map)
}

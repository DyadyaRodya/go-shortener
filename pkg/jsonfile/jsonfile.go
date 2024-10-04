package jsonfile

import (
	"encoding/json"
	"os"
)

func ReadFileToAny(Filename string, s any) error {
	file, err := os.OpenFile(Filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	return json.NewDecoder(file).Decode(s)
}

func WriteAnyToFile(Filename string, s any) error {
	file, err := os.OpenFile(Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil
	}

	defer file.Close()

	return json.NewEncoder(file).Encode(s)
}

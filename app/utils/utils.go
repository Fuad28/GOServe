package utils

import (
	"encoding/json"
	"io"
	"os"
)

func LoadFile[T any](filePath string, loadInto *T) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, loadInto)

	if err != nil {
		return err
	}

	return nil

}

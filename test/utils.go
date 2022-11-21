package test

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func FromTestData(path string, v any) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

    err = json.Unmarshal(bytes, &v)
	return err
}
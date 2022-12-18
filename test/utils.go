package test

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// UnmarshalFixture unmarshals a JSON file into a struct
func UnmarshalFixture(path string, v any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

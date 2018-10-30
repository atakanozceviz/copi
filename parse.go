package copi

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ParseIgnore file
func ParseIgnore(pth string) (map[string]struct{}, error) {
	conf := make(map[string]struct{})
	if pth == "" {
		return conf, nil
	}

	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// ParseTransform file
func ParseTransform(pth string) (map[string]string, error) {
	conf := make(map[string]string)
	if pth == "" {
		return conf, nil
	}

	f, err := os.Open(pth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

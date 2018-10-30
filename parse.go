package copi

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

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

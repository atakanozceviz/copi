package copi

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func parseSettings(pth string) (map[string]struct{}, error) {
	conf := make(map[string]struct{})

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
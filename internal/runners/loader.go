// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

type LoadOptions struct {
	Include []string
	Exclude []string
}

// loadAssets receives the Request argument and loads the assets from disk.
func loadAssets(in Request) (map[string]interface{}, error) {
	logger.Trace()

	path, err := importGetPathFromRequest(in)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var data = map[string]interface{}{}

	for _, ele := range files {
		var res interface{}
		if err := importLoadFromDisk(filepath.Join(path, ele.Name()), &res); err != nil {
			return nil, err
		}
		data[ele.Name()] = res
	}

	return data, nil
}

func loadStringAssets(in Request, options LoadOptions) (map[string]interface{}, error) {
	logger.Trace()

	path, err := importGetPathFromRequest(in)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var data = make(map[string]interface{})

	for _, ele := range files {
		res, err := utils.ReadStringFromFile(filepath.Join(path, ele.Name()))
		if err != nil {
			return nil, err
		}
		data[ele.Name()] = res
	}

	return data, nil

}

// loadUnmarshalAsset unmarshals the interface in argument to ptr
func loadUnmarshalAsset(in interface{}, ptr any) error {
	logger.Trace()

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &ptr); err != nil {
		return err
	}

	return nil
}

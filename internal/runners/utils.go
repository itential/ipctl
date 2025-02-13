// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
)

func GetProfile(name string, cfg *config.Config) (*config.Profile, error) {
	logger.Trace()

	active, err := cfg.ActiveProfile()

	profile, err := cfg.GetProfile(name)
	if err != nil {
		return nil, err
	}

	if active == profile {
		return nil, errors.New("source and destination servers are the same")
	}

	return profile, nil
}

func NewClient(name string, cfg *config.Config) (client.Client, context.CancelFunc, error) {
	logger.Trace()

	profile, err := cfg.GetProfile(name)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(profile.Timeout)*time.Second,
	)

	return client.New(ctx, profile), cancel, nil
}

func NormalizePath(in Request) (string, error) {
	logger.Trace()

	path := in.Args[0]

	if !utils.PathExists(path) {
		return "", errors.New("path does not exist")
	}

	return path, nil
}

// toMap will take any object and attempt to convert it into a map.   This
// function is primarily used to convert a struct into a map structure.  The
// function accepts a single argument `in` which is the struct to convert into
// a map.
func toMap(in any) (map[string]interface{}, error) {
	logger.Trace()

	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return res, nil
}

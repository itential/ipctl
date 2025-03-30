// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
)

// normalizeFilename will take a string argument for a filename and normalize
// it by replacing characters that will otherwise cause problems.
func normalizeFilename(s string) string {
	logger.Trace()
	return strings.Replace(s, "/", "_", -1)
}

// toMap takes a single required argument `in` and marshals it to a map using
// json.Marshal and json>Unmarshal.  This function will return the map or an
// error if one occurring during marshaling
func toMap(in any) (map[string]interface{}, error) {
	logger.Trace()

	var m map[string]interface{}
	if err := utils.ToMap(in, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func toArrayOfMaps(in any) ([]map[string]interface{}, error) {
	logger.Trace()

	var m []map[string]interface{}
	if err := utils.ToMap(in, &m); err != nil {
		return nil, err
	}

	return m, nil
}

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

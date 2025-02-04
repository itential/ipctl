// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"context"
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

/*
func NewServiceClient(name string, cfg *config.Config, f services.NewClientFunc) (services.Service, context.CancelFunc, error) {
	logger.Trace()

	c, cancel, err := NewClient(name, cfg)
	if err != nil {
		return nil, cancel, err
	}

	svc := f(c)

	return svc, cancel, nil
}
*/

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

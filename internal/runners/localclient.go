// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
)

type LocalClientRunner struct {
	config *config.Config
	client client.Client
}

func NewLocalClientRunner(client client.Client, cfg *config.Config) LocalClientRunner {
	return LocalClientRunner{
		config: cfg,
		client: client,
	}
}

func (r *LocalClientRunner) ShowConfig(in Request) (*Response, error) {
	logger.Trace()

	b, err := json.MarshalIndent(r.client.Profile(), "", "    ")
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   string(b),
		Object: r.client,
	}, nil
}

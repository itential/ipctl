// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"encoding/json"

	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
)

type LocalClientRunner struct {
	BaseRunner
	client client.Client
}

func NewLocalClientRunner(client client.Client, cfg *config.Config) LocalClientRunner {
	return LocalClientRunner{
		BaseRunner: NewBaseRunner(client, cfg),
		client:     client,
	}
}

func (r *LocalClientRunner) ShowConfig(in Request) (*Response, error) {
	logging.Trace()

	profile, err := r.config.ActiveProfile()
	if err != nil {
		logging.Warn("failed to load active profile, using defaults")
	}

	var port int

	if profile.UseTLS {
		port = 443
	} else if !profile.UseTLS {
		port = 80
	} else {
		port = profile.Port
	}

	config := map[string]interface{}{
		"host":          profile.Host,
		"port":          port,
		"use_tls":       profile.UseTLS,
		"verify":        profile.Verify,
		"username":      profile.Username,
		"password":      "********",
		"client_id":     profile.ClientID,
		"client_secret": "********",
	}

	b, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return nil, err
	}

	return &Response{
		Text:   string(b),
		Object: r.client,
	}, nil
}

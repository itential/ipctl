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

	return NewResponse(
		string(b),
		WithJson(r.client),
	), nil
}

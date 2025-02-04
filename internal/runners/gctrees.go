// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/services"
)

type GoldenConfigRunner struct {
	config  *config.Config
	service *services.GoldenConfigService
}

func NewGoldenConfigRunner(client client.Client, cfg *config.Config) *GoldenConfigRunner {
	return &GoldenConfigRunner{
		config:  cfg,
		service: services.NewGoldenConfigService(client),
	}
}

// Get implements the `get golden-config <name>` command
func (r *GoldenConfigRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	trees, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME"}
	for _, ele := range trees {
		display = append(display, ele.Name)
	}

	return NewResponse(
		strings.Join(display, "\n"),
		WithJson(trees),
	), nil
}

func (r *GoldenConfigRunner) Describe(in Request) (*Response, error) {
	return NotImplemented(in)
}

// Create implements the `create golden-config <name> <type>` command
func (r *GoldenConfigRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]
	deviceType := in.Args[1]

	gc, err := r.service.Create(services.GoldenConfigTree{
		Name:       name,
		DeviceType: deviceType,
	})
	if err != nil {
		return nil, err
	}

	return NewResponse(
		"Successfully create new golden configuration",
		WithJson(gc),
	), nil
}

// Delete implemetns the `delete golden-config <name>` command
func (r *GoldenConfigRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	treeId, err := r.getTreeIdFromName(name)
	if err != nil {
		return nil, err
	}

	if err := r.service.Delete(treeId); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully deleted Golden Config tree `%s`", name),
	), nil
}

func (r *GoldenConfigRunner) Clear(in Request) (*Response, error) {
	return NotImplemented(in)
}

func (r *GoldenConfigRunner) Copy(in Request) (*Response, error) {
	return NotImplemented(in)
}

func (r *GoldenConfigRunner) Import(in Request) (*Response, error) {
	return NotImplemented(in)
}

// Export is the implementation of the command `export golden-config <name>`
func (r *GoldenConfigRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetExportCommon
	utils.LoadObject(in.Common, &common)

	if common.Path != "" {
		utils.EnsurePathExists(common.Path)
	}

	name := in.Args[0]

	trees, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var id string
	for _, ele := range trees {
		if ele.Name == name {
			id = ele.Id
		}
	}

	if id == "" {
		return nil, errors.New(fmt.Sprintf("Unable to find tree with name %s", name))
	}

	tree, err := r.service.Export(id)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.gctree.json", name)
	if err := utils.WriteJsonToDisk(tree, fn, common.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported golden config tree `%s`", tree.Name),
	), nil
}

func (r *GoldenConfigRunner) getTreeIdFromName(name string) (string, error) {
	trees, err := r.service.GetAll()
	if err != nil {
		return "", err
	}

	for _, ele := range trees {
		if ele.Name == name {
			return ele.Id, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Unable to find tree with name %s", name))
}

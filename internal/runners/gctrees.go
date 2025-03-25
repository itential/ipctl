// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itential/ipctl/internal/flags"
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

/*
******************************************************************************
Reader interface
******************************************************************************
*/

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
		WithObject(trees),
	), nil
}

func (r *GoldenConfigRunner) Describe(in Request) (*Response, error) {
	return NotImplemented(in)
}

/*
******************************************************************************
Writer interface
******************************************************************************
*/

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
		WithObject(gc),
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

/*
******************************************************************************
Importer interface
******************************************************************************
*/

// Import implements the `import gctree ...` command
func (r *GoldenConfigRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var gctree services.GoldenConfigTree

	if err := importUnmarshalFromRequest(in, &gctree); err != nil {
		return nil, err
	}

	if err := r.importTree(gctree, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported gctree `%s`", gctree.Name),
		WithObject(gctree),
	), nil
}

/*
******************************************************************************
Exporter interface
******************************************************************************
*/

// Export implements the `export gctree ...` command
func (r *GoldenConfigRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	gctree, err := r.service.GetByName(name)
	if err != nil {
		return nil, err
	}

	res, err := r.service.Export(gctree.Id)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.gctree.json", res.Name)

	if err := exportAssetFromRequest(in, res, fn); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported gctree `%s`", gctree.Name),
		WithObject(gctree),
	), nil
}

/*
*******************************************************************************
Private functions
*******************************************************************************
*/

func (r *GoldenConfigRunner) importTree(in services.GoldenConfigTree, replace bool) error {
	logger.Trace()

	res, err := r.service.GetByName(in.Name)
	if err == nil {
		if res != nil {
			if replace {
				if err := r.service.Delete(res.Id); err != nil {
					return err
				}
			} else {
				return errors.New(
					fmt.Sprintf("gctree with name `%s` already exists, use `--replace` to overwrite it", in.Name),
				)
			}
		} else {
			return err
		}
	}

	return r.service.Import(in)
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

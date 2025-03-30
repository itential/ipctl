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
	"github.com/itential/ipctl/pkg/validators"
)

type ProfileRunner struct {
	config  *config.Config
	service *services.ProfileService
	client  client.Client
}

func NewProfileRunner(client client.Client, cfg *config.Config) *ProfileRunner {
	return &ProfileRunner{
		config:  cfg,
		service: services.NewProfileService(client),
		client:  client,
	}
}

//////////////////////////////////////////////////////////////////////////////
// Reader Interface
//

// Get implements the `get profiles` command
func (r *ProfileRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	profiles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	return &Response{
		Keys:   []string{"name", "description"},
		Object: profiles,
	}, nil
}

// Describe implements the `describe profile <name>` command
func (r *ProfileRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	profile, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return &Response{
		Object: profile,
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Writer Interface
//

// Create is the implementation of the command `ccreate profile <name>`
func (r *ProfileRunner) Create(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.ProfileCreateOptions
	utils.LoadObject(in.Options, &options)

	doc := services.NewProfile(name, options.Description)

	if _, err := r.service.Create(doc); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully created profile `%s`", name),
	}, nil
}

// Delete is the implementation of the command `delete profile <name>`
func (r *ProfileRunner) Delete(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	if err := r.service.Delete(name); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully deleted profile `%s`", name),
	}, nil
}

// Clear is the implementation of the command `clear profiles`
func (r *ProfileRunner) Clear(in Request) (*Response, error) {
	logger.Trace()

	cnt := 0

	profiles, err := r.service.GetAll()

	if err != nil {
		return nil, err
	}

	for _, ele := range profiles {
		r.service.Delete(ele.Id)
		cnt++
	}

	return &Response{
		Text: fmt.Sprintf("Deleted %v profile(s)", cnt),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Copier Interface
//

// Copy implements the `copy profile <name> <dst>` command
func (r *ProfileRunner) Copy(in Request) (*Response, error) {
	logger.Trace()

	res, err := Copy(CopyRequest{Request: in, Type: "profile"}, r)
	if err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully copied profile `%s` from `%s` to `%s`", res.Name, res.From, res.To),
	}, nil
}

func (r *ProfileRunner) CopyFrom(profile, name string) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	res, err := services.NewProfileService(client).Export(name)
	if err != nil {
		return nil, err
	}

	return *res, err
}

func (r *ProfileRunner) CopyTo(profile string, in any, replace bool) (any, error) {
	logger.Trace()

	client, cancel, err := NewClient(profile, r.config)
	if err != nil {
		return nil, err
	}
	defer cancel()

	svc := services.NewProfileService(client)

	id := in.(services.Profile).Id

	validator := validators.NewProfileValidator(r.client)

	if err := validator.CanImport(in.(services.Profile)); err != nil {
		if err.Code == validators.ValidatorAssetExists && replace {
			if err := svc.Delete(id); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return svc.Import(in.(services.Profile))
}

//////////////////////////////////////////////////////////////////////////////
// Importer Interface
//

// Import implements the command `import profile <path>`
func (r *ProfileRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	var profile services.Profile

	if err := importUnmarshalFromRequest(in, &profile); err != nil {
		return nil, err
	}

	if err := r.importProfile(profile, common.Replace); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully imported profile `%s`", profile.Id),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Exporter Interface
//

// Export is the implementation of the command `export profile <name>`
func (r *ProfileRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	profile, err := r.service.Export(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.profile.json", name)

	if err := exportAssetFromRequest(in, profile, fn); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Successfully exported profile `%s`", profile.Id),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Dumper Interface
//

// Dump implements the `dump prorfiles` command
func (r *ProfileRunner) Dump(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	var assets = map[string]interface{}{}
	for _, ele := range res {
		key := fmt.Sprintf("%s.profile.json", ele.Id)
		assets[key] = ele
	}

	if err := dumpAssets(in, assets); err != nil {
		return nil, err
	}

	return &Response{
		Text: fmt.Sprintf("Dumped %v profile(s)", len(assets)),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Loader Interface
//

// Load implements the `load profiles ...` command
func (r *ProfileRunner) Load(in Request) (*Response, error) {
	logger.Trace()

	elements, err := loadAssets(in)
	if err != nil {
		return nil, err
	}

	var loaded int
	var skipped int

	var output []string

	for fn, ele := range elements {
		var profile services.Profile

		if err := loadUnmarshalAsset(ele, &profile); err != nil {
			output = append(output, fmt.Sprintf("Failed to load profile from `%s`, skipping", fn))
			skipped++
			//return nil, err
		} else {
			if _, err := r.service.Import(profile); err != nil {
				if !strings.HasSuffix(err.Error(), "already exists!\"") {
					return nil, err
				}
				output = append(output, fmt.Sprintf("Skipping `%s`, profile `%s` already exists", fn, profile.Id))
				skipped++
			} else {
				output = append(output, fmt.Sprintf("Loaded profile `%s` successfully from `%s`", profile.Id, fn))
				loaded++
			}
		}
	}

	output = append(output, fmt.Sprintf(
		"\nSuccessfully loaded %v and skipped %v files from `%s`", loaded, skipped, in.Args[0],
	))

	return &Response{
		Text: strings.Join(output, "\n"),
	}, nil
}

//////////////////////////////////////////////////////////////////////////////
// Private functions
//

func (r *ProfileRunner) importProfile(in services.Profile, replace bool) error {
	logger.Trace()

	profile, err := r.service.Get(in.Id)
	if err != nil {
		return err
	}

	if profile != nil && replace {
		if err := r.service.Delete(profile.Id); err != nil {
			return err
		}
	} else {
		logger.Error(err, "")
		return errors.New(fmt.Sprintf("profile `%s` already exists", profile.Id))
	}

	if _, err := r.service.Import(in); err != nil {
		return err
	}

	return nil
}

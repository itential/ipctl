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

type ProfileRunner struct {
	config  *config.Config
	service *services.ProfileService
}

func NewProfileRunner(client client.Client, cfg *config.Config) *ProfileRunner {
	return &ProfileRunner{
		config:  cfg,
		service: services.NewProfileService(client),
	}
}

// Get implements the `get profiles` command
func (r *ProfileRunner) Get(in Request) (*Response, error) {
	logger.Trace()

	profiles, err := r.service.GetAll()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME\tDESCRIPTION"}
	for _, ele := range profiles {
		lines := []string{ele.Id, ele.Description}
		display = append(display, strings.Join(lines, "\t"))
	}

	return NewResponse(
		"",
		WithTable(display),
		WithJson(profiles),
	), nil
}

// Describe implements the `describe profile <name>` command
func (r *ProfileRunner) Describe(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	profile, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		"",
		WithJson(profile),
	), nil
}

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

	return NewResponse(
		fmt.Sprintf("Successfully created profile `%s`", name),
	), nil
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

	return NewResponse(fmt.Sprintf("Deleted %v profile(s)", cnt)), nil
}

// Copy implements the `copy profile <name> <dst>` command
func (r *ProfileRunner) Copy(in Request) (*Response, error) {
	logger.Trace()
	return NotImplemented(in)
}

// Import implements the command `import profile <path>`
func (r *ProfileRunner) Import(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetImportCommon
	utils.LoadObject(in.Common, &common)

	path, err := NormalizePath(in)
	if err != nil {
		return nil, err
	}

	data, err := utils.ReadFromFile(path)
	if err != nil {
		return nil, err
	}

	var profile services.Profile
	utils.UnmarshalData(data, &profile)

	imported, err := r.service.Import(profile)
	if err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully imported profile `%s`", imported.Id),
	), nil

}

// Export is the implementation of the command `export profile <name>`
func (r *ProfileRunner) Export(in Request) (*Response, error) {
	logger.Trace()

	name := in.Args[0]

	var options flags.AssetExportCommon
	utils.LoadObject(in.Common, &options)

	//profile, err := r.service.Export(name)
	profile, err := r.service.Get(name)
	if err != nil {
		return nil, err
	}

	fn := fmt.Sprintf("%s.profile.json", strings.Replace(profile.Id, "/", "_", 1))

	if err := utils.WriteJsonToDisk(profile, fn, options.Path); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully exported profile `%s` to `%s`", name, fn),
	), nil
}

// Pull implements the command `pull profile <repo>`
func (r *ProfileRunner) Pull(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetPullCommon
	utils.LoadObject(in.Common, &common)

	pull := PullAction{
		Name:     in.Args[1],
		Filename: in.Args[0],
		Config:   r.config,
		Options:  common,
	}

	data, err := pull.Do()
	if err != nil {
		return nil, err
	}

	var profile services.Profile
	utils.UnmarshalData(data, &profile)

	if err := r.importProfile(profile, common.Replace); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pulled profile `%s`", profile.Id),
	), nil
}

// Push implements the command `push profile <repo>`
func (r *ProfileRunner) Push(in Request) (*Response, error) {
	logger.Trace()

	var common flags.AssetPushCommon
	utils.LoadObject(in.Common, &common)

	//res, err := r.service.Export(in.Args[0])
	res, err := r.service.Get(in.Args[0])
	if err != nil {
		return nil, err
	}

	push := PushAction{
		Name:     in.Args[1],
		Filename: fmt.Sprintf("%s.profile.json", in.Args[0]),
		Options:  common,
		Config:   r.config,
		Data:     res,
	}

	if err := push.Do(); err != nil {
		return nil, err
	}

	return NewResponse(
		fmt.Sprintf("Successfully pushed profile `%s` to `%s`", in.Args[0], in.Args[1]),
	), nil
}

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
		return errors.New(fmt.Sprintf("profile `%s` already exists", profile.Id))
	}

	if _, err := r.service.Import(in); err != nil {
		return err
	}

	return nil
}

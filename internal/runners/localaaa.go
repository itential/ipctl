// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/terminal"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/localaaa"
	"github.com/itential/ipctl/pkg/logger"
)

type LocalAAARunner struct {
	config  *config.Config
	service localaaa.LocalAAAService
}

func NewLocalAAARunner(client client.Client, cfg *config.Config) LocalAAARunner {
	activeProfile, err := cfg.ActiveProfile()
	if err != nil {
		logger.Fatal(err, "")
	}

	mongoUrl := activeProfile.MongoUrl

	// NOTE (privateip) Disabling this block for the moment as it causes a
	// severe slow down running the app.  Need to re-think this
	/*
		if mongoUrl == "" {
			p, err := services.NewProfileService(client).GetActiveProfile()
			if err != nil {
				logger.Fatal(err, "failed to find the active profile")
			}

			adapters := p.Adapter["adapters"].([]interface{})

			for _, ele := range adapters {
				item := ele.(map[string]interface{})
				if item["type"].(string) == "local_aaa" {
					props := item["properties"].(map[string]interface{})
					db := props["database"].(map[string]interface{})
					mongoUrl = db["url"].(string)
				}
			}
		}
	*/

	runner := LocalAAARunner{config: cfg}

	if mongoUrl != "" {
		logger.Info("mongo url is %s", mongoUrl)
		runner.service = localaaa.NewLocalAAAService(mongoUrl)
	}

	return runner
}

func (r *LocalAAARunner) GetGroups(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetGroups()
	if err != nil {
		return nil, err
	}

	display := []string{"NAME"}
	for _, ele := range res {
		lines := []string{ele.Name}
		display = append(display, strings.Join(lines, "\t"))

	}

	return NewResponse(
		"",
		WithTable(display),
		WithObject(res),
	), nil
}

func (r *LocalAAARunner) CreateGroup(in Request) (*Response, error) {
	logger.Trace()

	grp := localaaa.NewGroup(in.Args[0])

	if err := r.service.CreateGroup(grp); err != nil {
		return nil, err
	}

	return NewResponse(
		"Successfully created new group",
	), nil

}

func (r *LocalAAARunner) DeleteGroup(in Request) (*Response, error) {
	logger.Trace()

	if err := r.service.DeleteGroup(in.Args[0]); err != nil {
		return nil, err
	}

	return NewResponse(
		"Successfully deleted group",
	), nil
}

func (r *LocalAAARunner) GetAccounts(in Request) (*Response, error) {
	logger.Trace()

	res, err := r.service.GetAccounts()
	if err != nil {
		return nil, err
	}

	display := []string{"USERNAME"}
	for _, ele := range res {
		lines := []string{ele.Username}
		display = append(display, strings.Join(lines, "\t"))

	}

	return NewResponse(
		"",
		WithTable(display),
		WithObject(res),
	), nil
}

func (r *LocalAAARunner) CreateAccount(in Request) (*Response, error) {
	logger.Trace()

	var options *flags.LocalAAAOptions
	utils.LoadObject(in.Options, &options)

	pw := terminal.Password()

	user := localaaa.NewAccount(in.Args[0], pw)
	user.Groups = options.Groups

	if err := r.service.CreateAccount(user); err != nil {
		return nil, err
	}

	return NewResponse(
		"Successfully created new user",
	), nil
}

func (r *LocalAAARunner) DeleteAccount(in Request) (*Response, error) {
	logger.Trace()

	if err := r.service.DeleteAccount(in.Args[0]); err != nil {
		return nil, err
	}

	return NewResponse(
		"Successfully deleted user",
	), nil
}

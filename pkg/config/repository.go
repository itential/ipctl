// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/go-homedir"
)

type Repository struct {
	Url            string `json:"url"`
	PrivateKey     string `json:"private_key"`
	PrivateKeyFile string `json:"private_key_file"`
	Reference      string `json:"reference"`
	Name           string `json:"name"`
	Email          string `json:"email"`
}

func getRepositoryFields() []string {
	fields := reflect.TypeOf((*Repository)(nil)).Elem()
	var properties = []string{}

	for i := 0; i < fields.NumField(); i++ {
		f := fields.Field(i)
		properties = append(properties, f.Tag.Get("json"))
	}

	return properties
}

func loadRepository(values, overrides map[string]interface{}) *Repository {
	r := &Repository{}

	for _, ele := range getRepositoryFields() {
		var value any

		// First check if the value for the key exists in overrides; if it
		// does, use it.  If it doesn't, fallback to checking if it wwas
		// defined in the profile (values)
		if val, ok := overrides[ele]; ok {
			value = val
		} else {
			value = values[ele]
		}

		var v string

		if value != nil {
			v = value.(string)
		}

		switch ele {
		case "url":
			r.Url = v
		case "private_key":
			r.PrivateKey = v
		case "private_key_file":
			privateKeyFile, _ := homedir.Expand(v)
			r.PrivateKeyFile = privateKeyFile
		case "reference":
			r.Reference = v
		case "name":
			r.Name = v
		case "email":
			r.Email = v
		}
	}

	return r
}

func (ac *Config) GetRepository(name string) (*Repository, error) {
	if v, exists := ac.repositories[name]; exists {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("repository `%s` does not exist", name))
}

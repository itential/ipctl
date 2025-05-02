// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/pflag"
)

const (
	defaultHost         = "localhost"
	defaultPort         = 0
	defaultUseTLS       = true
	defaultVerify       = true
	defaultUsername     = "admin@pronghorn"
	defaultPassword     = "admin"
	defaultClientID     = ""
	defaultClientSecret = ""
	defaultMongoUrl     = ""
	defaultTimeout      = 0
)

type Profile struct {
	Host string `json:"host"`
	Port int    `json:"port"`

	UseTLS bool `json:"use_tls"`
	Verify bool `json:"verify"`

	Username string `json:"username"`
	Password string `json:"password"`

	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	MongoUrl string `json:"mongo_url"`

	Timeout int `json:"timeout"`
}

func getProfileFields() []string {
	fields := reflect.TypeOf((*Profile)(nil)).Elem()
	var properties = []string{}

	for i := 0; i < fields.NumField(); i++ {
		f := fields.Field(i)
		properties = append(properties, f.Tag.Get("json"))
	}

	return properties
}

func loadProfile(values, defaults, overrides map[string]interface{}) *Profile {
	p := &Profile{}

	for _, ele := range getProfileFields() {
		var value any
		var exists bool

		// First check if the value for the key exists in overrides; if it
		// does, use it.  If it doesn't, fallback to checking if it wwas
		// defined in the profile (values).  If it was use it, if not, fall
		// back to defaults
		if val, ok := overrides[ele]; ok {
			value = val
		} else {
			value, exists = values[ele]
			if !exists {
				value = defaults[ele]
			}
		}

		var v string

		if value != nil {
			switch value.(type) {
			case int:
				v = strconv.Itoa(value.(int))
			case string:
				v = value.(string)
			case bool:
				if value.(bool) {
					v = "true"
				} else {
					v = "false"
				}
			default:
				panic(value)
			}
		}

		switch ele {
		case "host":
			if v == "" {
				v = defaultHost
			}
			p.Host = v
		case "port":
			if v != "" {
				port, err := strconv.Atoi(v)
				if err != nil {
					handleError("port must be valid integer", nil)
				}
				p.Port = port
			}
		case "use_tls":
			if v != "" {
				p.UseTLS = v == "true"
			} else {
				p.UseTLS = defaultUseTLS
			}
		case "verify":
			if v != "" {
				p.Verify = v == "true"
			} else {
				p.Verify = defaultVerify
			}
		case "username":
			if v == "" {
				v = defaultUsername
			}
			p.Username = v
		case "password":
			if v == "" {
				v = defaultPassword
			}
			p.Password = v
		case "client_id":
			if v == "" {
				v = defaultClientID
			}
			p.ClientID = v
		case "client_secret":
			if v == "" {
				v = defaultClientSecret
			}
			p.ClientSecret = v
		case "mongo_url":
			if v == "" {
				v = defaultMongoUrl
			}
			p.MongoUrl = v
		case "timeout":
			if v != "" {
				timeout, err := strconv.Atoi(v)
				if err != nil {
					handleError("timeout must be a valid integer", err)
				}
				p.Timeout = timeout
			} else {
				p.Timeout = defaultTimeout
			}
		}
	}

	return p
}

func getProfileFromFlag() string {
	var profile string

	profileFlag := pflag.NewFlagSet("profileFlag", pflag.ContinueOnError)
	profileFlag.StringVar(&profile, "profile", "", "Connection profile")
	profileFlag.ParseErrorsWhitelist.UnknownFlags = true
	profileFlag.Usage = func() {}

	// Non Nil ignore due to result empty slice
	if err := profileFlag.Parse(os.Args[1:]); err != nil && err != pflag.ErrHelp {
		handleError("failed to parse profile command line argument", err)
	}

	return profile
}

func (ac *Config) GetProfile(name string) (*Profile, error) {
	if cfg, exists := ac.profiles[name]; exists {
		return cfg, nil
	}
	p := loadProfile(
		map[string]interface{}{},
		map[string]interface{}{},
		map[string]interface{}{},
	)
	return p, fmt.Errorf("profile `%s` not found, using defaults", name)
}

func (ac *Config) ActiveProfile() (*Profile, error) {
	var profileName string
	if ac.profileName == "" {
		profileName = "default"
	} else {
		profileName = ac.profileName
	}
	return ac.GetProfile(profileName)
}

// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var profileFields = []string{
	"host",
	"port",
	"use_tls",
	"verify",
	"username",
	"password",
	"client_id",
	"client_secret",
	"timeout",
	"output",
	"verbose",
	"pager",
}

func TestGetProfileFields(t *testing.T) {
	for _, ele := range getProfileFields() {
		assert.Contains(t, profileFields, ele, "unknown field detected in profile")

	}
}

func TestLoadProfileWithValues(t *testing.T) {
	values := map[string]interface{}{
		"host":          "test",
		"port":          1000,
		"use_tls":       false,
		"verify":        true,
		"username":      "test",
		"password":      "test",
		"client_id":     "test",
		"client_secret": "test",
		"timeout":       1000,
	}

	p := loadProfile(values, map[string]interface{}{}, map[string]interface{}{})

	for key, value := range values {
		switch key {
		case "host":
			assert.Equal(t, p.Host, value)
		case "port":
			assert.Equal(t, p.Port, value)
		case "use_tls":
			assert.Equal(t, p.UseTLS, value)
		case "verify":
			assert.Equal(t, p.Verify, value)
		case "username":
			assert.Equal(t, p.Username, value)
		case "password":
			assert.Equal(t, p.Password, value)
		case "client_id":
			assert.Equal(t, p.ClientID, value)
		case "client_secret":
			assert.Equal(t, p.ClientSecret, value)
		case "timeout":
			assert.Equal(t, p.Timeout, value)
		}
	}
}

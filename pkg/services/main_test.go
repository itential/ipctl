// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

const (
	fixtureRoot = "testdata"
)

var (
	fixtureSuites = []string{"2023.2"}
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	zerolog.GlobalLevel()
	code := m.Run()
	os.Exit(code)
}

func fixtureDataToMap(data string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil, err
	}
	return m, nil
}

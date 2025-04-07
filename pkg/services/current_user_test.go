// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"path/filepath"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	getCurrentUserSuccess = "whoami.success.json"
)

func TestGetCurrentUser(t *testing.T) {
	testClient := testlib.Setup()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, getCurrentUserSuccess),
		)

		testlib.AddGetResponseToMux("/whoami", response, 0)

		res, err := GetCurrentUser(testClient)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "admin@pronghorn", res.Username)
	}
}

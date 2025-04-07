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
	userSettingsGetSuccess = "user/settings/get.success.json"
)

func setupUserSettingsService() *UserSettingsService {
	return NewUserSettingsService(
		testlib.Setup(),
	)
}

func TestUserSettingssGetAll(t *testing.T) {
	svc := setupUserSettingsService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, userSettingsGetSuccess),
		)

		data, err := fixtureDataToMap(response)
		if err != nil {
			t.FailNow()
		}

		id := data["_id"].(string)

		testlib.AddGetResponseToMux("/user/settings", response, 0)

		res, err := svc.Get()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id, res.Id)
	}
}

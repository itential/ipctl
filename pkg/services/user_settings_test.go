// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	userSettingsGetResponse = testlib.Fixture("testdata/user/settings/get.json")
)

func setupUserSettingsService() *UserSettingsService {
	return NewUserSettingsService(
		testlib.Setup(),
	)
}

func TestUserSettingssGetAll(t *testing.T) {
	svc := setupUserSettingsService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/user/settings", userSettingsGetResponse, 0)

	res, err := svc.Get()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.Id != "")
}

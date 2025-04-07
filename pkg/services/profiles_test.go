// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

const (
	profilesGetAllSuccess = "profiles/getall.success.json"
	profilesGetSuccess    = "profiles/get.success.json"
	profilesGetNotFound   = "profiles/get.notfound.json"
)

func setupProfileService() *ProfileService {
	return NewProfileService(
		testlib.Setup(),
	)
}

func TestProfileGetAll(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesGetAllSuccess),
		)
		testlib.AddGetResponseToMux("/profiles", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
	}
}

func TestProfileGetSuccess(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesGetSuccess),
		)
		testlib.AddGetResponseToMux("/profiles/{id}", response, 0)

		res, err := svc.Get("ID")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Profile)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Id != "")
	}
}

func TestProfileGetNotFound(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesGetNotFound),
		)

		testlib.AddGetResponseToMux("/profiles/{id}", response, 0)

		res, err := svc.Get("ID")

		assert.NotNil(t, err)
		assert.Nil(t, res)
	}
}

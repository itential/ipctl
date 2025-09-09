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
	profilesGetAllSuccess  = "profiles/getall.success.json"
	profilesGetSuccess     = "profiles/get.success.json"
	profilesGetNotFound    = "profiles/get.notfound.json"
	profilesCreateSuccess  = "profiles/create.json"
	profilesCreateExists   = "profiles/create.exists.json"
	profilesDeleteSuccess  = "profiles/delete.json"
	profilesDeleteNotFound = "profiles/delete.notfound.json"
	profilesImportSuccess  = "profiles/import.json"
	profilesImportExists   = "profiles/import.exists.json"
	profilesExportSuccess  = "profiles/export.json"
	profilesExportNotFound = "profiles/export.notfound.json"
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

func TestProfileGetActiveProfile(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesGetAllSuccess),
		)
		testlib.AddGetResponseToMux("/profiles", response, 0)

		res, err := svc.GetActiveProfile()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "profile1", res.Id)
	}
}

func TestProfileGetActiveProfileNotFound(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for range fixtureSuites {
		response := `{"results": [{"metadata": {"isActive": false, "activeSync": true}, "profile": {"id": "profile1"}}], "total": 1}`
		testlib.AddGetResponseToMux("/profiles", response, 0)

		res, err := svc.GetActiveProfile()

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "failed to find the active profile", err.Error())
	}
}

func TestProfileCreateSuccess(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesCreateSuccess),
		)
		testlib.AddPostResponseToMux("/profiles", response, 200)

		profile := NewProfile("test", "test description")
		res, err := svc.Create(profile)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.Id)
	}
}

func TestProfileDeleteSuccess(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesDeleteSuccess),
		)
		testlib.AddDeleteResponseToMux("/profiles/{name}", response, 0)

		err := svc.Delete("test")

		assert.Nil(t, err)
	}
}

func TestProfileDeleteNotFound(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesDeleteNotFound),
		)
		testlib.AddDeleteErrorToMux("/profiles/{name}", response, 404)

		err := svc.Delete("nonexistent")

		assert.NotNil(t, err)
	}
}

func TestProfileImportSuccess(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesImportSuccess),
		)
		testlib.AddPostResponseToMux("/profiles/import", response, 200)

		profile := NewProfile("test", "test description")
		res, err := svc.Import(profile)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.Id)
	}
}

func TestProfileExportSuccess(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesExportSuccess),
		)
		testlib.AddGetResponseToMux("/profiles/{name}", response, 0)

		res, err := svc.Export("test")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.Id)
	}
}

func TestProfileExportNotFound(t *testing.T) {
	svc := setupProfileService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, profilesExportNotFound),
		)
		testlib.AddGetErrorToMux("/profiles/{name}", response, 404)

		res, err := svc.Export("nonexistent")

		assert.NotNil(t, err)
		assert.Nil(t, res)
	}
}

// TODO: Add TestProfileActivateSuccess when testlib supports PUT requests

func TestNewProfile(t *testing.T) {
	name := "test-profile"
	desc := "test description"

	profile := NewProfile(name, desc)

	assert.Equal(t, name, profile.Id)
	assert.Equal(t, desc, profile.Description)
}

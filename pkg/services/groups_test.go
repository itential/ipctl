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

var (
	groupsGetSuccess    = "authorization/groups/get.success.json"
	groupsGetAllSuccess = "authorization/groups/getall.success.json"
)

func setupGroupService() *GroupService {
	return NewGroupService(
		testlib.Setup(),
	)
}

func TestGroupGetAll(t *testing.T) {
	svc := setupGroupService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, groupsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/authorization/groups", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
	}
}

func TestGroupGetAllError(t *testing.T) {
	svc := setupGroupService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authroization/groups", "", 0)

	res, err := svc.GetAll()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestGroupGet(t *testing.T) {
	svc := setupGroupService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, groupsGetSuccess),
		)
		testlib.AddGetResponseToMux("/authorization/groups/{id}", response, 0)

		res, err := svc.Get("ID")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Group)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Id != "")
	}
}

func TestGroupGetError(t *testing.T) {
	svc := setupGroupService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authorization/groups", "", 0)

	res, err := svc.Get("TEST")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestGroupGetByName(t *testing.T) {
	svc := setupGroupService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, groupsGetAllSuccess),
		)
		testlib.AddGetResponseToMux("/authorization/groups", response, 0)

		res, err := svc.GetByName("pronghorn_admin")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Group)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Id != "")
		assert.True(t, res.Name == "pronghorn_admin")
	}
}

func TestGroupGetByNameError(t *testing.T) {
	svc := setupGroupService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authorization/groups", "", 0)

	res, err := svc.GetByName("TEST")

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, reflect.TypeOf((*Group)(nil)), reflect.TypeOf(res))
}

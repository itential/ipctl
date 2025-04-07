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
	accountsGetSuccess    = "authorization/accounts/get.success.json"
	accountsGetAllSuccess = "authorization/accounts/getall.success.json"
)

func setupAccountService() *AccountService {
	return NewAccountService(
		testlib.Setup(),
	)
}

func TestAccountGetAll(t *testing.T) {
	svc := setupAccountService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, accountsGetAllSuccess),
		)
		testlib.AddGetResponseToMux("/authorization/accounts", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(res))
	}
}

func TestAccountGetAllError(t *testing.T) {
	svc := setupAccountService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authorization/accounts", "", 0)

	res, err := svc.GetAll()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAccountGet(t *testing.T) {
	svc := setupAccountService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, accountsGetSuccess),
		)

		testlib.AddGetResponseToMux("/authorization/accounts/{id}", response, 0)

		res, err := svc.Get("ID")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Account)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Id != "")
	}
}

func TestAccountGetError(t *testing.T) {
	svc := setupAccountService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authorization/accounts", "", 0)

	res, err := svc.Get("TEST")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAccountGetByName(t *testing.T) {
	svc := setupAccountService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, accountsGetAllSuccess),
		)
		testlib.AddGetResponseToMux("/authorization/accounts", response, 0)

		res, err := svc.GetByName("admin@pronghorn")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Account)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Id != "")
		assert.True(t, res.Username == "admin@pronghorn")

	}
}

func TestAccountGetByNameError(t *testing.T) {
	svc := setupAccountService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authorization/accounts", "", 0)

	res, err := svc.GetByName("TEST")

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, reflect.TypeOf((*Account)(nil)), reflect.TypeOf(res))
}

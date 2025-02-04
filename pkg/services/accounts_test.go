// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"reflect"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	accountsGetResponse    = testlib.Fixture("testdata/authorization/accounts/get.json")
	accountsGetAllResponse = testlib.Fixture("testdata/authorization/accounts/getall.json")
)

func setupAccountService() *AccountService {
	return NewAccountService(
		testlib.Setup(),
	)
}

func TestAccountGetAll(t *testing.T) {
	svc := setupAccountService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/authorization/accounts", accountsGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
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

	testlib.AddGetResponseToMux("/authorization/accounts/{id}", accountsGetResponse, 0)

	res, err := svc.Get("ID")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Account)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Id != "")
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

	testlib.AddGetResponseToMux("/authorization/accounts", accountsGetAllResponse, 0)

	res, err := svc.GetByName("admin@pronghorn")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Account)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Id != "")
	assert.True(t, res.Username == "admin@pronghorn")
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

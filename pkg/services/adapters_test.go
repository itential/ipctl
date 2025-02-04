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
	adaptersGetResponse    = testlib.Fixture("testdata/adapters/get.json")
	adaptersGetAllResponse = testlib.Fixture("testdata/adapters/getall.json")
)

func setupAdapterService() *AdapterService {
	return NewAdapterService(
		testlib.Setup(),
	)
}

func TestAdaptersGetAll(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/adapters", adaptersGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func TestAdapterGet(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/adapters/{name}", adaptersGetResponse, 0)

	res, err := svc.Get("local_aaa")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Adapter)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Name == "local_aaa")
}

func TestAdapterGetError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/adapters/{name}", "", 0)

	res, err := svc.Get("TEST")

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, reflect.TypeOf((*Adapter)(nil)), reflect.TypeOf(res))
}

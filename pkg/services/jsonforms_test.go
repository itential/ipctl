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
	jsonFormsGetAllSuccess = "json-forms/forms/getall.success.json"
	jsonFormsGetSuccess    = "json-forms/forms/get.success.json"
)

func setupJsonFormService() *JsonFormService {
	return NewJsonFormService(
		testlib.Setup(),
	)
}

func TestJsonFormGetAll(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/json-forms/forms", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(res))
	}
}

func TestJsonFormGet(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsGetSuccess),
		)

		testlib.AddGetResponseToMux("/json-forms/forms/test", response, 0)

		res, err := svc.Get("test")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*JsonForm)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Id != "")
		assert.True(t, res.Name == "test")
	}
}

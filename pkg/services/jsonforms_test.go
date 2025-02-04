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
	jsonFormsGetAllResponse = testlib.Fixture("testdata/json-forms/forms/getall.json")
	jsonFormsGetResponse    = testlib.Fixture("testdata/json-forms/forms/get.json")
)

func setupJsonFormService() *JsonFormService {
	return NewJsonFormService(
		testlib.Setup(),
	)
}

func TestJsonFormGetAll(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/json-forms/forms", jsonFormsGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
}

func TestJsonFormGet(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/json-forms/forms/test", jsonFormsGetResponse, 0)

	res, err := svc.Get("test")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*JsonForm)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Id != "")
	assert.True(t, res.Name == "test")
}

// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	transformationsGetAllResponse         = testlib.Fixture("testdata/transformations/getall.json")
	transformationsGetResponse            = testlib.Fixture("testdata/transformations/get.json")
	transformationsGetNotFoundResponse    = testlib.Fixture("testdata/transformations/get.notfound.json")
	transformationsCreateResponse         = testlib.Fixture("testdata/transformations/create.json")
	transformationsCreateErrorResponse    = testlib.Fixture("testdata/transformations/create.error.json")
	transformationsDeleteResponse         = testlib.Fixture("testdata/transformations/delete.json")
	transformationsDeleteNotFoundResponse = testlib.Fixture("testdata/transformations/delete.notfound.json")
	transformationsImportResponse         = testlib.Fixture("testdata/transformations/import.json")
)

func setupTransformationService() *TransformationService {
	return NewTransformationService(
		testlib.Setup(),
	)
}

func TestTransformationGetAll(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/transformations", transformationsGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func TestTransformationGet(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/transformations/679629d71190db2bc5752df2", transformationsGetResponse, 0)

	res, err := svc.Get("679629d71190db2bc5752df2")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Transformation)(nil)), reflect.TypeOf(res))
	assert.Equal(t, "679629d71190db2bc5752df2", res.Id)
}

func TestTransformationGetByName(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/transformations", transformationsGetAllResponse, 0)

	res, err := svc.GetByName("ipctl")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Transformation)(nil)), reflect.TypeOf(res))
	assert.Equal(t, "679629d71190db2bc5752df2", res.Id)

}

func TestTransformationGetByNameNotFound(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/transformations", transformationsGetNotFoundResponse, 0)

	res, err := svc.GetByName("test")

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, err.Error(), "transformation not found")
}

func TestTransformationCreate(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddPostResponseToMux("/transformations", transformationsCreateResponse, http.StatusOK)

	res, err := svc.Create(NewTransformation("ipctl", ""))

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Transformation)(nil)), reflect.TypeOf(res))
	assert.Equal(t, "679629d71190db2bc5752df2", res.Id)
}

func TestTransformationCreateDuplicate(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/transformations", transformationsCreateErrorResponse, 0)

	res, err := svc.Create(NewTransformation("ipctl", ""))

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestTransformationDelete(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddDeleteResponseToMux("/transformations/679629d71190db2bc5752df2", transformationsDeleteResponse, http.StatusNoContent)

	err := svc.Delete("679629d71190db2bc5752df2")

	assert.Nil(t, err)
}

func TestTransformationDeleteNotFound(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddDeleteErrorToMux("/transformations/test", transformationsDeleteNotFoundResponse, 0)

	err := svc.Delete("test")

	assert.NotNil(t, err)
}

func TestTransformationImport(t *testing.T) {
	svc := setupTransformationService()
	defer testlib.Teardown()

	testlib.AddPostResponseToMux("/transformations/import", transformationsImportResponse, http.StatusOK)

	res, err := svc.Import(NewTransformation("ipctl", ""))

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Transformation)(nil)), reflect.TypeOf(res))
	assert.Equal(t, "679635291190db2bc5752df3", res.Id)
}

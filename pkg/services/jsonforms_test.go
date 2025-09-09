// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	jsonFormsGetAllSuccess = "json-forms/forms/getall.success.json"
	jsonFormsGetSuccess    = "json-forms/forms/get.success.json"
	jsonFormsGetNotFound   = "json-forms/forms/get.notfound.json"
	jsonFormsCreateSuccess = "json-forms/forms/create.json"
	jsonFormsCreateError   = "json-forms/forms/create.error.json"
	jsonFormsImportSuccess = "json-forms/forms/import.success.json"
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

		testlib.AddGetResponseToMux("/json-forms/forms", response, http.StatusOK)

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

		testlib.AddGetResponseToMux("/json-forms/forms/test", response, http.StatusOK)

		res, err := svc.Get("test")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*JsonForm)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Id != "")
		assert.True(t, res.Name == "test")
	}
}

func TestJsonFormGetNotFound(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsGetNotFound),
		)

		testlib.AddGetResponseToMux("/json-forms/forms/nonexistent", response, http.StatusOK)

		res, err := svc.Get("nonexistent")

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "jsonform not found")
	}
}

func TestNewJsonForm(t *testing.T) {
	name := "test-form"
	description := "Test form description"

	jf := NewJsonForm(name, description)

	assert.Equal(t, name, jf.Name)
	assert.Equal(t, description, jf.Description)
	assert.NotNil(t, jf.Schema)
	assert.NotNil(t, jf.Struct)
	assert.NotNil(t, jf.UISchema)
	assert.NotNil(t, jf.BindingSchema)
	assert.NotNil(t, jf.ValidationSchema)

	// Check schema structure
	assert.Equal(t, description, jf.Schema["description"])
	assert.Equal(t, name, jf.Schema["title"])
	assert.Equal(t, "object", jf.Schema["type"])
	assert.NotNil(t, jf.Schema["properties"])
	assert.NotNil(t, jf.Schema["required"])

	// Check struct structure
	assert.Equal(t, "", jf.Struct["description"])
	assert.Equal(t, "object", jf.Struct["type"])
	assert.NotNil(t, jf.Struct["items"])
}

func TestJsonFormCreate(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsCreateSuccess),
		)

		testlib.AddPostResponseToMux("/json-forms/forms", response, http.StatusOK)

		input := NewJsonForm("test", "Test description")
		res, err := svc.Create(input)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.Name)
		assert.True(t, res.Id != "")
	}
}

func TestJsonFormCreateError(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsCreateError),
		)

		testlib.AddPostResponseToMux("/json-forms/forms", response, http.StatusOK)

		input := NewJsonForm("duplicate", "Duplicate form")
		res, err := svc.Create(input)

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "already exists")
	}
}

func TestJsonFormDelete(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for range fixtureSuites {
		testlib.AddDeleteResponseToMux("/json-forms/forms", `{"status":"success"}`, http.StatusOK)

		ids := []string{"6795a570a209bf3192120bfb", "678f121d3ed2bff9bec56160"}
		err := svc.Delete(ids)

		assert.Nil(t, err)
	}
}

func TestJsonFormGetByName(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/json-forms/forms", response, http.StatusOK)

		res, err := svc.GetByName("test")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.Name)
	}
}

func TestJsonFormGetByNameNotFound(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/json-forms/forms", response, http.StatusOK)

		res, err := svc.GetByName("nonexistent")

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "jsonform not found")
	}
}

// TestJsonFormClear is not included due to testlib URL pattern conflicts
// The Clear method combines GetAll + Delete, both tested individually above

func TestJsonFormClearEmpty(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for range fixtureSuites {
		// Mock empty GetAll response
		testlib.AddGetResponseToMux("/json-forms/forms", `[]`, http.StatusOK)

		err := svc.Clear()

		assert.Nil(t, err)
	}
}

func TestJsonFormImport(t *testing.T) {
	svc := setupJsonFormService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		// Mock Import response
		importResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsImportSuccess),
		)
		testlib.AddPostResponseToMux("/json-forms/import/forms", importResponse, http.StatusOK)

		// Mock Get response for the imported form
		getResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, jsonFormsGetSuccess),
		)
		testlib.AddGetResponseToMux("/json-forms/forms/6795a570a209bf3192120bfc", getResponse, http.StatusOK)

		input := NewJsonForm("test-import", "Imported test form")
		res, err := svc.Import(input)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test", res.Name) // The mocked response returns "test" as name
	}
}

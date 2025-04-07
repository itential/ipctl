// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	workflowsGetSuccess     = "automation-studio/workflows/get.success.json"
	workflowsGetAllSuccess  = "automation-studio/workflows/getall.success.json"
	workflowsDeleteSuccess  = "automation-studio/workflows/delete.success.json"
	workflowsDeleteNotFound = "automation-studio/workflows/delete.notfound.json"
	workflowsCreateSuccess  = "automation-studio/automations/create.success.json"
	workflowsExportSuccess  = "workflow_builder/export/export.success.json"
	workflowsImportSuccess  = "automation-studio/automations/import.success.json"
	workflowsImportError    = "automation-studio/automations/import.error.json"
)

func setupWorkflowService() *WorkflowService {
	return NewWorkflowService(
		testlib.Setup(),
	)
}

func TestWorkflowGetAll(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/automation-studio/workflows", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 5, len(res))
	}
}

func TestWorkflowGet(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsGetSuccess),
		)

		testlib.AddGetResponseToMux("/automation-studio/workflows", response, 0)

		data, err := fixtureDataToMap(response)
		if err != nil {
			t.FailNow()
		}

		items := data["items"].([]interface{})

		name := items[0].(map[string]interface{})["name"].(string)
		id := items[0].(map[string]interface{})["_id"].(string)

		res, err := svc.Get(name)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
		assert.Equal(t, id, res.Id)
		assert.Equal(t, name, res.Name)
	}

}

func TestWorkflowGetError(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/automation-studio/workflows", "", 0)

	res, err := svc.Get("test")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestWorkflowCreate(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsCreateSuccess),
		)

		testlib.AddPostResponseToMux("/automation-studio/automations", response, http.StatusOK)

		data, err := fixtureDataToMap(response)
		if err != nil {
			t.FailNow()
		}

		name := data["created"].(map[string]interface{})["name"].(string)
		id := data["created"].(map[string]interface{})["_id"].(string)

		doc := NewWorkflow(name)

		res, err := svc.Create(doc)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
		assert.Equal(t, name, res.Name)
		assert.Equal(t, id, res.Id)
	}
}

func TestWorkflowCreateError(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/automation-studio/workflows", "", 0)

	doc := NewWorkflow("TEST")

	res, err := svc.Create(doc)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestWorkflowDelete(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsDeleteSuccess),
		)

		data, err := fixtureDataToMap(response)
		if err != nil {
			t.FailNow()
		}

		name := data["value"].(map[string]interface{})["name"].(string)

		testlib.AddDeleteResponseToMux(
			fmt.Sprintf("/workflow_builder/workflows/delete/%s", name), response, http.StatusOK,
		)

		err = svc.Delete(name)

		assert.Nil(t, err)
	}
}

func TestWorkflowDeleteNotFound(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsDeleteNotFound),
		)

		testlib.AddDeleteErrorToMux("/workflow_builder/workflows/delete/test", response, 0)

		err := svc.Delete("test")

		assert.NotNil(t, err)
	}
}

func TestWorkflowExport(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsExportSuccess),
		)

		data, err := fixtureDataToMap(response)
		if err != nil {
			t.FailNow()
		}

		name := data["name"].(string)

		testlib.AddPostErrorToMux("/workflow_builder/export", response, http.StatusOK)

		res, err := svc.Export(name)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
		assert.Equal(t, name, res.Name)
	}
}

func TestWorkflowImport(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsImportSuccess),
		)

		data, err := fixtureDataToMap(response)
		if err != nil {
			t.FailNow()
		}

		name := data["name"].(string)

		testlib.AddPostResponseToMux("/automation-studio/automations/import", response, http.StatusOK)

		doc := NewWorkflow(name)

		res, err := svc.Import(doc)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
		assert.Equal(t, name, res.Name)
	}
}

func TestWorkflowImportError(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, workflowsImportError),
		)

		testlib.AddPostResponseToMux("/automation-studio/automations/import", response, http.StatusInternalServerError)

		doc := NewWorkflow("test")

		res, err := svc.Import(doc)

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
	}
}

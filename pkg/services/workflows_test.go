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
	workflowsGetResponse           = testlib.Fixture("testdata/automation-studio/workflows/get.json")
	workflowsGetAllResponse        = testlib.Fixture("testdata/automation-studio/workflows/getall.json")
	workflowsDeleteSuccessResponse = testlib.Fixture("testdata/automation-studio/workflows/delete/success.json")
	workflowsDeleteErrorResponse   = testlib.Fixture("testdata/automation-studio/workflows/delete/error.json")
	workflowsCreateResponse        = testlib.Fixture("testdata/automation-studio/automations/create.json")
	workflowsExportResponse        = testlib.Fixture("testdata/workflow_builder/export/post.json")
	workflowsImportSuccessResponse = testlib.Fixture("testdata/automation-studio/automations/import.success.json")
	workflowsImportErrorResponse   = testlib.Fixture("testdata/automation-studio/automations/import.error.json")
)

func setupWorkflowService() *WorkflowService {
	return NewWorkflowService(
		testlib.Setup(),
	)
}

func TestWorkflowGetAll(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/automation-studio/workflows", workflowsGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 5, len(res))
}

func TestWorkflowGet(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/automation-studio/workflows", workflowsGetResponse, 0)

	res, err := svc.Get("UI-Test-1")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Id != "")
}

func TestWorkflowGetError(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/automation-studio/workflows", "", 0)

	res, err := svc.Get("TEST")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestWorkflowCreate(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddPostResponseToMux("/automation-studio/automations", workflowsCreateResponse, http.StatusOK)

	doc := NewWorkflow("UI-Test-1")

	res, err := svc.Create(doc)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Id != "")
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

	testlib.AddDeleteResponseToMux("/workflow_builder/workflows/delete/test", workflowsDeleteSuccessResponse, http.StatusOK)

	err := svc.Delete("test")

	assert.Nil(t, err)
}

func TestWorkflowDeleteError(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddDeleteErrorToMux("/workflow_builder/workflows/delete/test", workflowsDeleteErrorResponse, 0)

	err := svc.Delete("test")

	assert.NotNil(t, err)
}

func TestWorkflowExport(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/workflow_builder/export", workflowsExportResponse, http.StatusOK)

	res, err := svc.Export("test")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Name == "test")
}

func TestWorkflowImport(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddPostResponseToMux("/automation-studio/automations/import", workflowsImportSuccessResponse, http.StatusOK)

	doc := NewWorkflow("test")

	res, err := svc.Import(doc)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
	assert.True(t, res.Name == "test")
}

func TestWorkflowImportError(t *testing.T) {
	svc := setupWorkflowService()
	defer testlib.Teardown()

	testlib.AddPostResponseToMux("/automation-studio/automations/import", workflowsImportErrorResponse, http.StatusInternalServerError)

	doc := NewWorkflow("test")

	res, err := svc.Import(doc)

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, reflect.TypeOf((*Workflow)(nil)), reflect.TypeOf(res))
}

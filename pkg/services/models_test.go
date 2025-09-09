// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	modelsGetAllSuccess    = "lifecycle-manager/getall.success.json"
	modelsGetSuccess       = "lifecycle-manager/get.success.json"
	modelsCreateSuccess    = "lifecycle-manager/create.success.json"
	modelsImportSuccess    = "lifecycle-manager/import.success.json"
	modelsExportSuccess    = "lifecycle-manager/export.success.json"
	modelsRunActionSuccess = "lifecycle-manager/runaction.success.json"
)

func setupModelService() *ModelService {
	return NewModelService(
		testlib.Setup(),
	)
}

func TestNewModel(t *testing.T) {
	name := "test-model"
	desc := "Test model description"

	model := NewModel(name, desc)

	assert.Equal(t, name, model.Name)
	assert.Equal(t, desc, model.Description)
	assert.Empty(t, model.Id)
	assert.Nil(t, model.Schema)
	assert.Empty(t, model.Actions)
}

func TestNewModelService(t *testing.T) {
	client := testlib.Setup()
	svc := NewModelService(client)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.client)
}

func TestModelService_GetAll(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/lifecycle-manager/resources", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, "test-model-1", res[0].Name)
		assert.Equal(t, "test-model-2", res[1].Name)
	}
}

func TestModelService_Get(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	modelID := "64f1c2b8e4b0123456789abc"

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsGetSuccess),
		)

		testlib.AddGetResponseToMux("/lifecycle-manager/resources/"+modelID, response, 0)

		res, err := svc.Get(modelID)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, modelID, res.Id)
		assert.Equal(t, "test-model", res.Name)
		assert.Equal(t, "A test model for unit testing", res.Description)
		assert.NotNil(t, res.Schema)
		assert.Equal(t, 1, len(res.Actions))
	}
}

func TestModelService_GetByName(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	modelName := "test-model-1"

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/lifecycle-manager/resources", response, 0)

		res, err := svc.GetByName(modelName)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "64f1c2b8e4b0123456789abc", res.Id)
		assert.Equal(t, modelName, res.Name)
	}
}

func TestModelService_GetByName_NotFound(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	modelName := "nonexistent-model"

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/lifecycle-manager/resources", response, 0)

		res, err := svc.GetByName(modelName)

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "model not found", err.Error())
	}
}

func TestModelService_Create(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	model := Model{
		Name:        "new-model",
		Description: "A newly created model",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"field1": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsCreateSuccess),
		)

		testlib.AddPostResponseToMux("/lifecycle-manager/resources", response, http.StatusOK)

		res, err := svc.Create(model)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "64f1c2b8e4b0123456789new", res.Id)
		assert.Equal(t, "new-model", res.Name)
		assert.Equal(t, "A newly created model", res.Description)
	}
}

func TestModelService_Create_WithoutSchema(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	model := Model{
		Name:        "new-model",
		Description: "A newly created model",
	}

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsCreateSuccess),
		)

		testlib.AddPostResponseToMux("/lifecycle-manager/resources", response, http.StatusOK)

		res, err := svc.Create(model)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "new-model", res.Name)
	}
}

func TestModelService_Delete(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	modelID := "64f1c2b8e4b0123456789abc"

	testlib.AddDeleteResponseToMux("/lifecycle-manager/resources/"+modelID, "", http.StatusOK)

	err := svc.Delete(modelID, false)

	assert.Nil(t, err)
}

func TestModelService_Delete_WithInstances(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	modelID := "64f1c2b8e4b0123456789abc"

	testlib.AddDeleteResponseToMux("/lifecycle-manager/resources/"+modelID, "", http.StatusOK)

	err := svc.Delete(modelID, true)

	assert.Nil(t, err)
}

func TestModelService_RunAction(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	modelID := "64f1c2b8e4b0123456789abc"
	request := RunActionRequest{
		ActionId: "action1",
		Inputs: map[string]interface{}{
			"input1": "value1",
		},
		Name:         "test-action",
		Instance:     "instance123",
		InstanceName: "test-instance",
	}

	response := testlib.Fixture(
		filepath.Join(fixtureRoot, fixtureSuites[0], modelsRunActionSuccess),
	)

	testlib.AddPostResponseToMux("/lifecycle-manager/resources/"+modelID+"/run-action", response, http.StatusCreated)

	res, err := svc.RunAction(modelID, request)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "run123", res.Id)
	assert.Equal(t, modelID, res.ModelId)
	assert.Equal(t, "action1", res.ActionId)
	assert.Equal(t, "completed", res.Status)
	assert.Equal(t, "job123", res.JobId)
}

func TestModelService_Import(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	model := Model{
		Name:        "imported-model",
		Description: "An imported model",
		Schema: map[string]interface{}{
			"type": "object",
		},
	}

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsImportSuccess),
		)

		testlib.AddPostResponseToMux("/lifecycle-manager/resources/import", response, http.StatusOK)

		res, err := svc.Import(model)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "64f1c2b8e4b0123456789imp", res.Id)
		assert.Equal(t, "imported-model", res.Name)
		assert.Equal(t, "An imported model", res.Description)
	}
}

func TestModelService_Export(t *testing.T) {
	svc := setupModelService()
	defer testlib.Teardown()

	modelID := "64f1c2b8e4b0123456789exp"

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, modelsExportSuccess),
		)

		testlib.AddGetResponseToMux("/lifecycle-manager/resources/"+modelID+"/export", response, 0)

		res, err := svc.Export(modelID)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, modelID, res.Id)
		assert.Equal(t, "exported-model", res.Name)
		assert.Equal(t, "An exported model", res.Description)
		assert.Equal(t, 1, len(res.Actions))
	}
}

func TestRunActionRequest_Structure(t *testing.T) {
	req := RunActionRequest{
		ActionId: "test-action",
		Inputs: map[string]interface{}{
			"input1": "value1",
			"input2": 123,
		},
		Name:         "test-name",
		Instance:     "test-instance",
		InstanceName: "test-instance-name",
	}

	assert.Equal(t, "test-action", req.ActionId)
	assert.Equal(t, "value1", req.Inputs["input1"])
	assert.Equal(t, 123, req.Inputs["input2"])
	assert.Equal(t, "test-name", req.Name)
	assert.Equal(t, "test-instance", req.Instance)
	assert.Equal(t, "test-instance-name", req.InstanceName)
}

func TestRunActionResponse_Structure(t *testing.T) {
	resp := RunActionResponse{
		Id:         "response-id",
		ModelId:    "model-id",
		InstanceId: "instance-id",
		ActionId:   "action-id",
		Status:     "completed",
		Progress: map[string]interface{}{
			"percentage": 100,
		},
		Errors: []string{},
		InitialInstanceData: map[string]interface{}{
			"field1": "initial",
		},
		FinalInstanceData: map[string]interface{}{
			"field1": "final",
		},
	}

	assert.Equal(t, "response-id", resp.Id)
	assert.Equal(t, "model-id", resp.ModelId)
	assert.Equal(t, "completed", resp.Status)
	assert.Equal(t, 100, resp.Progress["percentage"])
	assert.Equal(t, 0, len(resp.Errors))
}

func TestModelAction_Structure(t *testing.T) {
	action := ModelAction{
		Id:              "action-id",
		Name:            "action-name",
		PreWorkflowJst:  "pre-workflow",
		PostWorkflowJst: "post-workflow",
		Workflow:        "workflow-id",
		Type:            "workflow",
	}

	assert.Equal(t, "action-id", action.Id)
	assert.Equal(t, "action-name", action.Name)
	assert.Equal(t, "pre-workflow", action.PreWorkflowJst)
	assert.Equal(t, "post-workflow", action.PostWorkflowJst)
	assert.Equal(t, "workflow-id", action.Workflow)
	assert.Equal(t, "workflow", action.Type)
}

func TestModel_Structure(t *testing.T) {
	model := Model{
		Id:          "model-id",
		Name:        "model-name",
		Description: "model-description",
		Schema: map[string]interface{}{
			"type": "object",
		},
		Actions: []ModelAction{
			{
				Id:   "action1",
				Name: "First Action",
			},
		},
		Created:       "2023-09-01T12:00:00.000Z",
		CreatedBy:     "user123",
		LastUpdated:   "2023-09-01T13:00:00.000Z",
		LastUpdatedBy: "user123",
	}

	assert.Equal(t, "model-id", model.Id)
	assert.Equal(t, "model-name", model.Name)
	assert.Equal(t, "model-description", model.Description)
	assert.Equal(t, "object", model.Schema["type"])
	assert.Equal(t, 1, len(model.Actions))
	assert.Equal(t, "action1", model.Actions[0].Id)
}

func TestModelOperation_Structure(t *testing.T) {
	op := ModelOperation{
		Message: "Operation successful",
		Data: []Model{
			{
				Id:   "model1",
				Name: "Model 1",
			},
			{
				Id:   "model2",
				Name: "Model 2",
			},
		},
		Metadata: Metadata{
			Total: 2,
		},
	}

	assert.Equal(t, "Operation successful", op.Message)
	assert.Equal(t, 2, len(op.Data))
	assert.Equal(t, "model1", op.Data[0].Id)
	assert.Equal(t, "Model 1", op.Data[0].Name)
	assert.Equal(t, 2, op.Metadata.Total)
}

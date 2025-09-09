// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	automationsGetAllSuccess = "automation-studio/automations/create.success.json"
	automationsGetSuccess    = "automation-studio/automations/create.success.json"
	automationsCreateSuccess = "automation-studio/automations/create.success.json"
	automationsImportSuccess = "automation-studio/automations/import.success.json"
	automationsImportError   = "automation-studio/automations/import.error.json"
)

func setupAutomationService() *AutomationService {
	return NewAutomationService(
		testlib.Setup(),
	)
}

func TestAutomationService_GetAll(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	mockResponse := `{"message":"Successfully retrieved automations","data":[{"_id":"test-id","name":"Test Automation","description":"Test Description","componentName":"test-component","componentType":"workflows","componentId":"component-id","gbac":{"read":[],"write":[]},"created":"2023-01-01T00:00:00Z","createdBy":"test-user","lastUpdated":"2023-01-01T00:00:00Z","lastUpdatedBy":"test-user"}],"metadata":{"total":1,"skip":0,"limit":100}}`

	testlib.AddGetResponseToMux("/operations-manager/automations", mockResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, len(res))
	if len(res) > 0 {
		assert.Equal(t, "test-id", res[0].Id)
		assert.Equal(t, "Test Automation", res[0].Name)
	}
}

func TestAutomationService_Get(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	mockResponse := `{"message":"Successfully retrieved automation","data":{"_id":"test-id","name":"Test Automation","description":"Test Description","componentName":"test-component","componentType":"workflows","componentId":"component-id","gbac":{"read":[],"write":[]},"created":"2023-01-01T00:00:00Z","createdBy":"test-user","lastUpdated":"2023-01-01T00:00:00Z","lastUpdatedBy":"test-user"}}`

	testlib.AddGetResponseToMux("/operations-manager/automations/test-id", mockResponse, 0)

	res, err := svc.Get("test-id")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "test-id", res.Id)
	assert.Equal(t, "Test Automation", res.Name)
}

func TestAutomationService_GetByName(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	mockResponse := `{"message":"Successfully retrieved automations","data":[{"_id":"test-id","name":"Test Automation","description":"Test Description","componentName":"test-component","componentType":"workflows","componentId":"component-id","gbac":{"read":[],"write":[]},"created":"2023-01-01T00:00:00Z","createdBy":"test-user","lastUpdated":"2023-01-01T00:00:00Z","lastUpdatedBy":"test-user"}],"metadata":{"total":1,"skip":0,"limit":100}}`

	testlib.AddGetResponseToMux("/operations-manager/automations", mockResponse, 0)

	res, err := svc.GetByName("Test Automation")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Test Automation", res.Name)
}

func TestAutomationService_Create(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	mockResponse := `{"message":"Successfully created automation","data":{"_id":"new-automation-id","name":"Test Automation","description":"Test Description","componentType":"workflows"}}`

	testlib.AddPostResponseToMux("/operations-manager/automations", mockResponse, 200)

	automation := NewAutomation("Test Automation", "Test Description")
	res, err := svc.Create(automation)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Test Automation", res.Name)
}

func TestAutomationService_Delete(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	testlib.AddDeleteResponseToMux("/operations-manager/automations/test-id", "", 0)

	err := svc.Delete("test-id")
	assert.Nil(t, err)
}

func TestAutomationService_Import_Success(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	mockResponse := `{"message":"Successfully imported automation","data":[{"success":true,"data":{"_id":"imported-id","name":"Test Import","description":"Test Import Description"}}]}`

	testlib.AddPutResponseToMux("/operations-manager/automations", mockResponse, 0)

	automation := NewAutomation("Test Import", "Test Import Description")
	res, err := svc.Import(automation)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestAutomationService_Export(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	mockResponse := `{"message":"Successfully exported automation","data":{"_id":"test-id","name":"Exported Automation","triggers":[]}}`

	testlib.AddGetResponseToMux("/operations-manager/automations/test-id/export", mockResponse, 0)

	res, err := svc.Export("test-id")

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestAutomationService_Clear(t *testing.T) {
	svc := setupAutomationService()
	defer testlib.Teardown()

	mockResponse := `{"message":"Successfully retrieved automations","data":[{"_id":"test-id","name":"Test Automation"}],"metadata":{"total":1}}`

	testlib.AddGetResponseToMux("/operations-manager/automations", mockResponse, 0)
	testlib.AddDeleteResponseToMux("/operations-manager/automations/test-id", "", 0)

	err := svc.Clear()
	assert.Nil(t, err)
}

func TestNewAutomation(t *testing.T) {
	automation := NewAutomation("Test", "Description")

	assert.Equal(t, "Test", automation.Name)
	assert.Equal(t, "Description", automation.Description)
	assert.Equal(t, "workflows", automation.ComponentType)
}

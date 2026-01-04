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
	integrationsGetSuccess    = "integrations/get.success.json"
	integrationsGetAllSuccess = "integrations/getall.success.json"
	integrationsCreateSuccess = "integrations/create.success.json"
)

func setupIntegrationService() *IntegrationService {
	return NewIntegrationService(
		testlib.Setup(),
	)
}

// TestNewIntegrationService tests the constructor
func TestNewIntegrationService(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewIntegrationService(client)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.client)
	assert.Equal(t, reflect.TypeOf((*IntegrationService)(nil)), reflect.TypeOf(svc))
}

// TestNewIntegration tests the integration factory function
func TestNewIntegration(t *testing.T) {
	name := "test-integration"
	integrationType := "http"

	integration := NewIntegration(name, integrationType)

	assert.Equal(t, name, integration.Name)
	assert.NotNil(t, integration.Properties)
	assert.Equal(t, name, integration.Properties["id"])
	assert.Equal(t, integrationType, integration.Properties["type"])

	// Test with different parameters
	integration2 := NewIntegration("db-integration", "database")
	assert.Equal(t, "db-integration", integration2.Name)
	assert.Equal(t, "db-integration", integration2.Properties["id"])
	assert.Equal(t, "database", integration2.Properties["type"])
}

// TestNewIntegrationEmptyParams tests NewIntegration with empty parameters
func TestNewIntegrationEmptyParams(t *testing.T) {
	integration := NewIntegration("", "")

	assert.Empty(t, integration.Name)
	assert.NotNil(t, integration.Properties)
	assert.Empty(t, integration.Properties["id"])
	assert.Empty(t, integration.Properties["type"])
}

// TestIntegrationServiceGetAll tests retrieving all integrations
func TestIntegrationServiceGetAll(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/integrations", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(res))

		// Verify first integration
		assert.Equal(t, "integration-1", res[0].Name)
		assert.Equal(t, "Adapter", res[0].Type)
		assert.Equal(t, "http-adapter", res[0].Model)
		assert.Equal(t, true, res[0].Virtual)
		assert.Equal(t, false, res[0].IsEncrypted)
		assert.Equal(t, false, res[0].Profiling)
		assert.NotNil(t, res[0].Properties)
		assert.Equal(t, "integration-1", res[0].Properties["id"])
		assert.Equal(t, "http", res[0].Properties["type"])

		// Verify second integration
		assert.Equal(t, "integration-2", res[1].Name)
		assert.Equal(t, "Adapter", res[1].Type)
		assert.Equal(t, "database-adapter", res[1].Model)
		assert.Equal(t, true, res[1].Virtual)
		assert.Equal(t, true, res[1].IsEncrypted)
		assert.Equal(t, true, res[1].Profiling)
	}
}

// TestIntegrationServiceGetAllError tests error handling for GetAll
func TestIntegrationServiceGetAllError(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/integrations", "", 0)

	res, err := svc.GetAll()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestIntegrationServiceGet tests retrieving a specific integration by name
func TestIntegrationServiceGet(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationsGetSuccess),
		)
		testlib.AddGetResponseToMux("/integrations/{name}", response, 0)

		res, err := svc.Get("test-integration")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Integration)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "test-integration", res.Name)
		assert.Equal(t, "Adapter", res.Type)
		assert.Equal(t, "test-model", res.Model)
		assert.Equal(t, true, res.Virtual)
		assert.Equal(t, false, res.IsEncrypted)
		assert.Equal(t, false, res.Profiling)

		// Verify properties
		assert.NotNil(t, res.Properties)
		assert.Equal(t, "test-integration", res.Properties["id"])
		assert.Equal(t, "http", res.Properties["type"])
		assert.Equal(t, "localhost", res.Properties["host"])
		assert.Equal(t, float64(8080), res.Properties["port"])

		// Verify logger properties
		assert.NotNil(t, res.LoggerProperties)
		assert.Equal(t, "info", res.LoggerProperties["level"])
		assert.Equal(t, true, res.LoggerProperties["enabled"])

		// Verify system properties
		assert.NotNil(t, res.SystemProperties)
		assert.Equal(t, float64(5000), res.SystemProperties["timeout"])
		assert.Equal(t, float64(3), res.SystemProperties["retries"])

		// Verify event deduplication
		assert.NotNil(t, res.EventDeduplciation)
		assert.Equal(t, false, res.EventDeduplciation["enabled"])
		assert.Equal(t, float64(60000), res.EventDeduplciation["window"])
	}
}

// TestIntegrationServiceGetError tests error handling for Get
func TestIntegrationServiceGetError(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/integrations/{name}", "", 0)

	res, err := svc.Get("non-existent")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestIntegrationServiceCreate tests creating a new integration
func TestIntegrationServiceCreate(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationsCreateSuccess),
		)
		testlib.AddPostResponseToMux("/integrations", response, http.StatusOK)

		inputIntegration := Integration{
			Name:  "new-integration",
			Model: "new-model",
			Properties: map[string]interface{}{
				"id":       "new-integration",
				"type":     "api",
				"endpoint": "https://api.newservice.com",
			},
			IsEncrypted: false,
			LoggerProperties: map[string]interface{}{
				"level":   "info",
				"enabled": true,
			},
			Profiling: false,
			SystemProperties: map[string]interface{}{
				"timeout": 30000,
				"retries": 5,
			},
			EventDeduplciation: map[string]interface{}{
				"enabled": true,
				"window":  120000,
			},
		}

		res, err := svc.Create(inputIntegration)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "new-integration", res.Name)
		assert.Equal(t, "Adapter", res.Type) // Should be set automatically
		assert.Equal(t, "new-model", res.Model)
		assert.Equal(t, true, res.Virtual) // Should be set automatically
		assert.Equal(t, false, res.IsEncrypted)
		assert.Equal(t, false, res.Profiling)

		// Verify properties
		assert.NotNil(t, res.Properties)
		assert.Equal(t, "new-integration", res.Properties["id"])
		assert.Equal(t, "api", res.Properties["type"])
		assert.Equal(t, "https://api.newservice.com", res.Properties["endpoint"])
	}
}

// TestIntegrationServiceCreateAutoFields tests that Create sets required fields
func TestIntegrationServiceCreateAutoFields(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationsCreateSuccess),
		)
		testlib.AddPostResponseToMux("/integrations", response, http.StatusOK)

		// Create integration with different Type and Virtual values
		inputIntegration := Integration{
			Name:    "test-integration",
			Type:    "SomeOtherType", // Should be overridden
			Virtual: false,           // Should be overridden
			Properties: map[string]interface{}{
				"id": "test-integration",
			},
		}

		res, err := svc.Create(inputIntegration)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		// The returned data should reflect the API response, not necessarily the overridden input
		assert.Equal(t, "Adapter", res.Type)
		assert.Equal(t, true, res.Virtual)
	}
}

// TestIntegrationServiceCreateError tests error handling for Create
func TestIntegrationServiceCreateError(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/integrations", "", http.StatusInternalServerError)

	inputIntegration := Integration{
		Name: "test-integration",
		Properties: map[string]interface{}{
			"id": "test-integration",
		},
	}

	res, err := svc.Create(inputIntegration)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestIntegrationServiceDelete tests deleting an integration
func TestIntegrationServiceDelete(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	testlib.AddDeleteResponseToMux("/integrations/{name}", "", http.StatusOK)

	err := svc.Delete("test-integration")

	assert.Nil(t, err)
}

// TestIntegrationServiceDeleteError tests error handling for Delete
func TestIntegrationServiceDeleteError(t *testing.T) {
	svc := setupIntegrationService()
	defer testlib.Teardown()

	testlib.AddDeleteErrorToMux("/integrations/{name}", "", http.StatusNotFound)

	err := svc.Delete("non-existent")

	assert.NotNil(t, err)
}

// TestIntegrationStructFields tests the Integration struct field mappings
func TestIntegrationStructFields(t *testing.T) {
	integration := Integration{
		Name:  "test-integration",
		Type:  "Adapter",
		Model: "test-model",
		Properties: map[string]interface{}{
			"id":   "test-integration",
			"type": "http",
		},
		IsEncrypted: true,
		LoggerProperties: map[string]interface{}{
			"level": "debug",
		},
		Virtual:   true,
		Profiling: true,
		SystemProperties: map[string]interface{}{
			"timeout": 5000,
		},
		EventDeduplciation: map[string]interface{}{
			"enabled": true,
		},
	}

	assert.Equal(t, "test-integration", integration.Name)
	assert.Equal(t, "Adapter", integration.Type)
	assert.Equal(t, "test-model", integration.Model)
	assert.Equal(t, true, integration.IsEncrypted)
	assert.Equal(t, true, integration.Virtual)
	assert.Equal(t, true, integration.Profiling)
	assert.NotNil(t, integration.Properties)
	assert.NotNil(t, integration.LoggerProperties)
	assert.NotNil(t, integration.SystemProperties)
	assert.NotNil(t, integration.EventDeduplciation)
	assert.Equal(t, "test-integration", integration.Properties["id"])
	assert.Equal(t, "debug", integration.LoggerProperties["level"])
	assert.Equal(t, 5000, integration.SystemProperties["timeout"])
	assert.Equal(t, true, integration.EventDeduplciation["enabled"])
}

// TestIntegrationStructEmptyFields tests Integration with empty/nil fields
func TestIntegrationStructEmptyFields(t *testing.T) {
	integration := Integration{}

	assert.Empty(t, integration.Name)
	assert.Empty(t, integration.Type)
	assert.Empty(t, integration.Model)
	assert.False(t, integration.IsEncrypted)
	assert.False(t, integration.Virtual)
	assert.False(t, integration.Profiling)
	assert.Nil(t, integration.Properties)
	assert.Nil(t, integration.LoggerProperties)
	assert.Nil(t, integration.SystemProperties)
	assert.Nil(t, integration.EventDeduplciation)
}

// TestIntegrationServiceClientType tests that the service uses the correct client type
func TestIntegrationServiceClientType(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewIntegrationService(client)

	// Verify the service was created successfully with embedded BaseService
	assert.NotNil(t, svc)
	assert.Equal(t, reflect.TypeOf((*IntegrationService)(nil)), reflect.TypeOf(svc))
}

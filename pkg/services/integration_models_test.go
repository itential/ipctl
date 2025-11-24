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
	integrationModelsGetSuccess        = "integration-models/get.success.json"
	integrationModelsGetAllSuccess     = "integration-models/getall.success.json"
	integrationModelsCreateSuccess     = "integration-models/create.success.json"
	integrationModelsGetCreatedSuccess = "integration-models/get-created.success.json"
	integrationModelsExportSuccess     = "integration-models/export.success.json"
)

func setupIntegrationModelService() *IntegrationModelService {
	return NewIntegrationModelService(
		testlib.Setup(),
	)
}

// TestNewIntegrationModelService tests the constructor
func TestNewIntegrationModelService(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewIntegrationModelService(client)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.client)
	assert.Equal(t, reflect.TypeOf((*IntegrationModelService)(nil)), reflect.TypeOf(svc))
}

// TestIntegrationModelServiceGetAll tests retrieving all integration models
func TestIntegrationModelServiceGetAll(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationModelsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/integration-models", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 3, len(res))

		// Verify first integration model
		assert.Equal(t, "http-adapter", res[0].Model)
		assert.Equal(t, "http-adapter-v1.0.0", res[0].VersionId)
		assert.Equal(t, "HTTP adapter for REST API integrations", res[0].Description)
		assert.Equal(t, "1.0.0", res[0].Properties.Version)
		assert.Equal(t, "https", res[0].Properties.Server.Protocol)
		assert.Equal(t, "example.com", res[0].Properties.Server.Host)
		assert.Equal(t, "/api", res[0].Properties.Server.BasePath)
		assert.Equal(t, true, res[0].Properties.Tls.Enabled)
		assert.Equal(t, false, res[0].Properties.Tls.RefjectUnauthorized)
		assert.NotNil(t, res[0].Properties.Authentication)

		// Verify second integration model
		assert.Equal(t, "database-adapter", res[1].Model)
		assert.Equal(t, "database-adapter-v2.1.0", res[1].VersionId)
		assert.Equal(t, "Database adapter for SQL database integrations", res[1].Description)
		assert.Equal(t, "2.1.0", res[1].Properties.Version)
		assert.Equal(t, "tcp", res[1].Properties.Server.Protocol)
		assert.Equal(t, false, res[1].Properties.Tls.Enabled)
		assert.Equal(t, true, res[1].Properties.Tls.RefjectUnauthorized)

		// Verify third integration model
		assert.Equal(t, "file-adapter", res[2].Model)
		assert.Equal(t, "file-adapter-v1.5.0", res[2].VersionId)
		assert.Equal(t, "1.5.0", res[2].Properties.Version)
		assert.Equal(t, "file", res[2].Properties.Server.Protocol)
		assert.Equal(t, "localhost", res[2].Properties.Server.Host)
	}
}

// TestIntegrationModelServiceGetAllError tests error handling for GetAll
func TestIntegrationModelServiceGetAllError(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/integration-models", "", 0)

	res, err := svc.GetAll()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestIntegrationModelServiceGet tests retrieving a specific integration model by name
func TestIntegrationModelServiceGet(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationModelsGetSuccess),
		)
		testlib.AddGetResponseToMux("/integration-models/{name}", response, 0)

		res, err := svc.Get("http-adapter")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*IntegrationModel)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "http-adapter", res.Model)
		assert.Equal(t, "http-adapter-v1.0.0", res.VersionId)
		assert.Equal(t, "HTTP adapter model for REST API integrations", res.Description)

		// Verify properties structure
		assert.Equal(t, "1.0.0", res.Properties.Version)
		assert.NotNil(t, res.Properties.Authentication)

		// Verify server properties
		assert.Equal(t, "https", res.Properties.Server.Protocol)
		assert.Equal(t, "api.example.com", res.Properties.Server.Host)
		assert.Equal(t, "/api/v1", res.Properties.Server.BasePath)

		// Verify TLS properties
		assert.Equal(t, true, res.Properties.Tls.Enabled)
		assert.Equal(t, false, res.Properties.Tls.RefjectUnauthorized)
	}
}

// TestIntegrationModelServiceGetError tests error handling for Get
func TestIntegrationModelServiceGetError(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/integration-models/{name}", "", 0)

	res, err := svc.Get("non-existent")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestIntegrationModelServiceCreate tests creating a new integration model
func TestIntegrationModelServiceCreate(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		// Mock the POST response
		createResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationModelsCreateSuccess),
		)
		testlib.AddPostResponseToMux("/integration-models", createResponse, http.StatusOK)

		// Mock the GET response for the created model
		getResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationModelsGetCreatedSuccess),
		)
		testlib.AddGetResponseToMux("/integration-models/{name}", getResponse, 0)

		inputModel := map[string]interface{}{
			"model":       "custom-adapter",
			"description": "Custom adapter for testing",
			"properties": map[string]interface{}{
				"authentication": map[string]interface{}{
					"type":  "token",
					"token": "test-token",
				},
				"server": map[string]interface{}{
					"protocol":  "https",
					"host":      "custom.example.com",
					"base_path": "/custom/api",
				},
				"tls": map[string]interface{}{
					"enabled":            true,
					"rejectUnauthroized": true,
				},
				"version": "1.0.0",
			},
		}

		res, err := svc.Create(inputModel)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "custom-adapter", res.Model)
		assert.Equal(t, "custom-adapter-v1.0.0", res.VersionId)
		assert.Equal(t, "Custom adapter model created for testing", res.Description)
		assert.Equal(t, "1.0.0", res.Properties.Version)
		assert.Equal(t, "https", res.Properties.Server.Protocol)
		assert.Equal(t, "custom.example.com", res.Properties.Server.Host)
		assert.Equal(t, true, res.Properties.Tls.Enabled)
		assert.Equal(t, true, res.Properties.Tls.RefjectUnauthorized)
	}
}

// TestIntegrationModelServiceCreateError tests error handling for Create
func TestIntegrationModelServiceCreateError(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/integration-models", "", http.StatusInternalServerError)

	inputModel := map[string]interface{}{
		"model": "test-adapter",
	}

	res, err := svc.Create(inputModel)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestIntegrationModelServiceCreateGetError tests Create when follow-up Get fails
func TestIntegrationModelServiceCreateGetError(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		// Mock successful POST response
		createResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationModelsCreateSuccess),
		)
		testlib.AddPostResponseToMux("/integration-models", createResponse, http.StatusOK)

		// Mock failed GET response
		testlib.AddGetErrorToMux("/integration-models/{name}", "", 0)

		inputModel := map[string]interface{}{
			"model": "test-adapter",
		}

		res, err := svc.Create(inputModel)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	}
}

// TestIntegrationModelServiceDelete tests deleting an integration model
func TestIntegrationModelServiceDelete(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	testlib.AddDeleteResponseToMux("/integration-models/{name}", "", http.StatusOK)

	err := svc.Delete("http-adapter")

	assert.Nil(t, err)
}

// TestIntegrationModelServiceDeleteError tests error handling for Delete
func TestIntegrationModelServiceDeleteError(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	testlib.AddDeleteErrorToMux("/integration-models/{name}", "", http.StatusNotFound)

	err := svc.Delete("non-existent")

	assert.NotNil(t, err)
}

// TestIntegrationModelServiceExport tests exporting an integration model
func TestIntegrationModelServiceExport(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, integrationModelsExportSuccess),
		)
		testlib.AddGetResponseToMux("/integration-models/{name}/export", response, 0)

		res, err := svc.Export("http-adapter")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "http-adapter", res["model"])
		assert.Equal(t, "1.0.0", res["version"])
		assert.NotNil(t, res["metadata"])
		assert.NotNil(t, res["definition"])
		assert.NotNil(t, res["schema"])

		// Verify metadata structure
		metadata, ok := res["metadata"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "2023-01-01T12:00:00Z", metadata["created"])
		assert.Equal(t, "system", metadata["author"])

		// Verify definition structure
		definition, ok := res["definition"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "HTTP Adapter", definition["name"])
		assert.Equal(t, "Comprehensive HTTP adapter for REST API integrations", definition["description"])
		assert.NotNil(t, definition["properties"])
	}
}

// TestIntegrationModelServiceExportError tests error handling for Export
func TestIntegrationModelServiceExportError(t *testing.T) {
	svc := setupIntegrationModelService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/integration-models/{name}/export", "", 0)

	res, err := svc.Export("non-existent")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestIntegrationModelTlsStruct tests the IntegrationModelTls struct
func TestIntegrationModelTlsStruct(t *testing.T) {
	tls := IntegrationModelTls{
		Enabled:             true,
		RefjectUnauthorized: false,
	}

	assert.Equal(t, true, tls.Enabled)
	assert.Equal(t, false, tls.RefjectUnauthorized)
}

// TestIntegrationModelServerStruct tests the IntegrationModelServer struct
func TestIntegrationModelServerStruct(t *testing.T) {
	server := IntegrationModelServer{
		Protocol: "https",
		Host:     "api.example.com",
		BasePath: "/api/v1",
	}

	assert.Equal(t, "https", server.Protocol)
	assert.Equal(t, "api.example.com", server.Host)
	assert.Equal(t, "/api/v1", server.BasePath)
}

// TestIntegrationModelPropertiesStruct tests the IntegrationModelProperties struct
func TestIntegrationModelPropertiesStruct(t *testing.T) {
	properties := IntegrationModelProperties{
		Authentication: map[string]interface{}{
			"type":     "basic",
			"username": "admin",
		},
		Server: IntegrationModelServer{
			Protocol: "https",
			Host:     "example.com",
			BasePath: "/api",
		},
		Tls: IntegrationModelTls{
			Enabled:             true,
			RefjectUnauthorized: false,
		},
		Version: "2.0.0",
	}

	assert.NotNil(t, properties.Authentication)
	assert.Equal(t, "basic", properties.Authentication["type"])
	assert.Equal(t, "https", properties.Server.Protocol)
	assert.Equal(t, true, properties.Tls.Enabled)
	assert.Equal(t, "2.0.0", properties.Version)
}

// TestIntegrationModelStruct tests the IntegrationModel struct
func TestIntegrationModelStruct(t *testing.T) {
	model := IntegrationModel{
		Model:       "test-adapter",
		VersionId:   "test-adapter-v1.0.0",
		Description: "Test adapter model",
		Properties: IntegrationModelProperties{
			Version: "1.0.0",
			Server: IntegrationModelServer{
				Protocol: "http",
				Host:     "localhost",
				BasePath: "/test",
			},
			Tls: IntegrationModelTls{
				Enabled:             false,
				RefjectUnauthorized: true,
			},
			Authentication: map[string]interface{}{
				"type": "none",
			},
		},
	}

	assert.Equal(t, "test-adapter", model.Model)
	assert.Equal(t, "test-adapter-v1.0.0", model.VersionId)
	assert.Equal(t, "Test adapter model", model.Description)
	assert.Equal(t, "1.0.0", model.Properties.Version)
	assert.Equal(t, "http", model.Properties.Server.Protocol)
	assert.Equal(t, false, model.Properties.Tls.Enabled)
	assert.NotNil(t, model.Properties.Authentication)
}

// TestIntegrationModelStructEmptyFields tests IntegrationModel with empty fields
func TestIntegrationModelStructEmptyFields(t *testing.T) {
	model := IntegrationModel{}

	assert.Empty(t, model.Model)
	assert.Empty(t, model.VersionId)
	assert.Empty(t, model.Description)
	assert.Empty(t, model.Properties.Version)
	assert.Empty(t, model.Properties.Server.Protocol)
	assert.Empty(t, model.Properties.Server.Host)
	assert.False(t, model.Properties.Tls.Enabled)
	assert.False(t, model.Properties.Tls.RefjectUnauthorized)
	assert.Nil(t, model.Properties.Authentication)
}

// TestIntegrationModelServiceClientType tests that the service uses the correct client type
func TestIntegrationModelServiceClientType(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewIntegrationModelService(client)

	assert.NotNil(t, svc.client)
	assert.Equal(t, reflect.TypeOf((*ServiceClient)(nil)), reflect.TypeOf(svc.client))
}

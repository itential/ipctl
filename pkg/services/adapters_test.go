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
	adaptersGetSuccess    = "adapters/get.success.json"
	adaptersGetAllSuccess = "adapters/getall.success.json"
)

func setupAdapterService() *AdapterService {
	return NewAdapterService(
		testlib.Setup(),
	)
}

func TestAdaptersGetAll(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, adaptersGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/adapters", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
	}
}

func TestAdapterGet(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, adaptersGetSuccess),
		)

		testlib.AddGetResponseToMux("/adapters/{name}", response, 0)

		res, err := svc.Get("local_aaa")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Adapter)(nil)), reflect.TypeOf(res))
		assert.True(t, res.Name == "local_aaa")
	}
}

func TestAdapterGetError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/adapters/{name}", "", 0)

	res, err := svc.Get("TEST")

	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.Equal(t, reflect.TypeOf((*Adapter)(nil)), reflect.TypeOf(res))
}

func TestAdapterCreate(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	// Mock response for successful create
	response := `{
		"status": "success",
		"message": "Adapter created successfully",
		"data": {
			"name": "test_adapter",
			"type": "Adapter",
			"model": "@itential/adapter-test",
			"properties": {
				"id": "test_adapter",
				"type": "test",
				"brokers": ["test"],
				"groups": [],
				"properties": {}
			},
			"isEncrypted": false,
			"loggerProps": {},
			"virtual": false
		}
	}`

	testlib.AddPostResponseToMux("/adapters", response, 200)

	inputAdapter := Adapter{
		Name: "test_adapter",
		Type: "Adapter",
		Model: "@itential/adapter-test",
		Properties: AdapterProperties{
			Id:         "test_adapter",
			Type:       "test",
			Brokers:    []string{"test"},
			Groups:     []any{},
			Properties: map[string]interface{}{},
		},
		IsEncrypted:      false,
		LoggerProperties: map[string]interface{}{},
		Virtual:          false,
	}

	res, err := svc.Create(inputAdapter)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "test_adapter", res.Name)
	assert.Equal(t, "Adapter", res.Type)
}

func TestAdapterCreateError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/adapters", "Internal server error", 500)

	inputAdapter := Adapter{
		Name: "test_adapter",
	}

	res, err := svc.Create(inputAdapter)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAdapterDelete(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddDeleteResponseToMux("/adapters/{name}", "", 200)

	err := svc.Delete("test_adapter")

	assert.Nil(t, err)
}

func TestAdapterDeleteError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddDeleteErrorToMux("/adapters/{name}", "Not found", 404)

	err := svc.Delete("nonexistent_adapter")

	assert.NotNil(t, err)
}

func TestAdapterImport(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	// Mock response for successful import
	response := `{
		"status": "success",
		"message": "Adapter imported successfully",
		"data": {
			"name": "imported_adapter",
			"type": "Adapter",
			"model": "@itential/adapter-imported",
			"properties": {
				"id": "imported_adapter",
				"type": "imported",
				"brokers": ["test"],
				"groups": [],
				"properties": {}
			},
			"isEncrypted": false,
			"loggerProps": {},
			"virtual": false
		}
	}`

	testlib.AddPostResponseToMux("/adapters/import", response, 200)

	inputAdapter := Adapter{
		Name: "imported_adapter",
		Type: "Adapter",
		Model: "@itential/adapter-imported",
	}

	res, err := svc.Import(inputAdapter)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "imported_adapter", res.Name)
}

func TestAdapterImportError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/adapters/import", "Bad request", 400)

	inputAdapter := Adapter{
		Name: "invalid_adapter",
	}

	res, err := svc.Import(inputAdapter)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAdapterUpdate(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	// Mock response for successful update
	response := `{
		"status": "success",
		"message": "Adapter updated successfully",
		"data": {
			"name": "updated_adapter",
			"type": "Adapter",
			"model": "@itential/adapter-updated",
			"properties": {
				"id": "updated_adapter",
				"type": "updated",
				"brokers": ["test"],
				"groups": [],
				"properties": {}
			},
			"isEncrypted": true,
			"loggerProps": {},
			"virtual": false
		}
	}`

	testlib.AddPutResponseToMux("/adapters/{name}", response, 200)

	inputAdapter := Adapter{
		Name: "updated_adapter",
		Type: "Adapter",
		Model: "@itential/adapter-updated",
		IsEncrypted: true,
	}

	res, err := svc.Update(inputAdapter)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "updated_adapter", res.Name)
	assert.Equal(t, true, res.IsEncrypted)
}

func TestAdapterUpdateError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddPutErrorToMux("/adapters/{name}", "Not found", 404)

	inputAdapter := Adapter{
		Name: "nonexistent_adapter",
	}

	res, err := svc.Update(inputAdapter)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAdapterExport(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, adaptersGetSuccess),
		)

		testlib.AddGetResponseToMux("/adapters/{name}", response, 0)

		res, err := svc.Export("local_aaa")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "local_aaa", res.Name)
	}
}

func TestAdapterExportError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/adapters/{name}", "Not found", 404)

	res, err := svc.Export("nonexistent_adapter")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestAdapterStart(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddPutResponseToMux("/adapters/{name}/start", "", 200)

	err := svc.Start("test_adapter")

	assert.Nil(t, err)
}

func TestAdapterStartError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddPutErrorToMux("/adapters/{name}/start", "Not found", 404)

	err := svc.Start("nonexistent_adapter")

	assert.NotNil(t, err)
}

func TestAdapterStop(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddPutResponseToMux("/adapters/{name}/stop", "", 200)

	err := svc.Stop("test_adapter")

	assert.Nil(t, err)
}

func TestAdapterStopError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	testlib.AddPutErrorToMux("/adapters/{name}/stop", "Not found", 404)

	err := svc.Stop("nonexistent_adapter")

	assert.NotNil(t, err)
}

func TestAdapterRestart(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	// Mock both stop and start operations
	testlib.AddPutResponseToMux("/adapters/{name}/stop", "", 200)
	testlib.AddPutResponseToMux("/adapters/{name}/start", "", 200)

	err := svc.Restart("test_adapter")

	assert.Nil(t, err)
}

func TestAdapterRestartStopError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	// Mock stop to fail
	testlib.AddPutErrorToMux("/adapters/{name}/stop", "Stop failed", 500)

	err := svc.Restart("test_adapter")

	assert.NotNil(t, err)
}

func TestAdapterRestartStartError(t *testing.T) {
	svc := setupAdapterService()
	defer testlib.Teardown()

	// Mock stop to succeed but start to fail
	testlib.AddPutResponseToMux("/adapters/{name}/stop", "", 200)
	testlib.AddPutErrorToMux("/adapters/{name}/start", "Start failed", 500)

	err := svc.Restart("test_adapter")

	assert.NotNil(t, err)
}

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
	projectsGetAllSuccess = "automation-studio/projects/getall_response.json"
	projectsGetSuccess    = "automation-studio/projects/get_response.json"
	projectsCreateSuccess = "automation-studio/projects/create_response.json"
	projectsExportSuccess = "automation-studio/projects/export_response.json"
)

func setupProjectService() *ProjectService {
	return NewProjectService(
		testlib.Setup(),
	)
}

func TestNewProjectService(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewProjectService(client)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.client)
	assert.Equal(t, reflect.TypeOf((*ProjectService)(nil)), reflect.TypeOf(svc))
}

func TestProjectsGetAll(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, projectsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/automation-studio/projects", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 7, len(res))
		assert.Equal(t, "Port/VLAN Configuration - IOS", res[0].Name)
		assert.Equal(t, "66aba9be41f8aad085ca0ee3", res[0].Id)
	}
}

func TestProjectsGetAllError(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/automation-studio/projects", "Internal server error", 500)

	res, err := svc.GetAll()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestProjectsGet(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, projectsGetSuccess),
		)

		testlib.AddGetResponseToMux("/automation-studio/projects/670d8dac113f9679380359de", response, 0)

		res, err := svc.Get("670d8dac113f9679380359de")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Project)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "Untitled-Project", res.Name)
		assert.Equal(t, "670d8dac113f9679380359de", res.Id)
		assert.Equal(t, 139, res.Iid)
	}
}

func TestProjectsGetError(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/automation-studio/projects/nonexistent", "Not found", 404)

	res, err := svc.Get("nonexistent")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestProjectsCreate(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	response := `{
		"message": "Project created successfully",
		"data": {
			"_id": "test-project-id",
			"name": "TestProject",
			"description": "",
			"components": [],
			"folders": [],
			"iid": 123,
			"componentIidIndex": 0,
			"created": "2024-01-01T00:00:00.000Z",
			"createdBy": {
				"_id": "user-id",
				"username": "testuser"
			},
			"lastUpdated": "2024-01-01T00:00:00.000Z",
			"lastUpdatedBy": {
				"_id": "user-id",
				"username": "testuser"
			}
		}
	}`

	testlib.AddPostResponseToMux("/automation-studio/projects", response, http.StatusOK)

	res, err := svc.Create("TestProject")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "TestProject", res.Name)
	assert.Equal(t, "test-project-id", res.Id)
}

func TestProjectsCreateError(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/automation-studio/projects", "Bad request", 400)

	res, err := svc.Create("TestProject")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestProjectsDelete(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddDeleteResponseToMux("/automation-studio/projects/test-project-id", "", 200)

	err := svc.Delete("test-project-id")

	assert.Nil(t, err)
}

func TestProjectsDeleteError(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddDeleteErrorToMux("/automation-studio/projects/nonexistent", "Not found", 404)

	err := svc.Delete("nonexistent")

	assert.NotNil(t, err)
}

func TestProjectsExport(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, projectsGetSuccess),
		)

		testlib.AddGetResponseToMux("/automation-studio/projects/670d8dac113f9679380359de/export", response, 0)

		res, err := svc.Export("670d8dac113f9679380359de")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Untitled-Project", res.Name)
		assert.Equal(t, "670d8dac113f9679380359de", res.Id)
	}
}

func TestProjectsExportError(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/automation-studio/projects/nonexistent/export", "Not found", 404)

	res, err := svc.Export("nonexistent")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestProjectImportMethod(t *testing.T) {
	project := Project{
		Id:                "test-id",
		Name:              "TestProject",
		BackgroundColor:   "#ffffff",
		Components:        []ProjectComponent{},
		Created:           "2024-01-01T00:00:00.000Z",
		CreatedBy:         "test-user",
		Description:       "Test description",
		Folders:           []ProjectFolder{},
		Iid:               123,
		ComponentIidIndex: 0, // Should be excluded
		LastUpdated:       "2024-01-01T00:00:00.000Z",
		LastUpdatedBy:     "test-user",
		Thumbnail:         "test-thumbnail",
		Members:           []ProjectMember{},      // Should be excluded
		AccessControl:     ProjectAccessControl{}, // Should be excluded
	}

	result := project.Import()

	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result["_id"])
	assert.Equal(t, "TestProject", result["name"])
	assert.Equal(t, "#ffffff", result["backgroundColor"])
	assert.Equal(t, "Test description", result["description"])
	assert.Equal(t, 123, result["iid"])
	assert.Equal(t, "test-thumbnail", result["thumbnail"])

	// These fields should not be present in import
	_, hasComponentIidIndex := result["componentIidIndex"]
	_, hasMembers := result["members"]
	_, hasAccessControl := result["accessControl"]

	assert.False(t, hasComponentIidIndex)
	assert.False(t, hasMembers)
	assert.False(t, hasAccessControl)
}

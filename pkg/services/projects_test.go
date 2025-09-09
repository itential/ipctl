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
	projectsGetAllSuccess     = "automation-studio/projects/getall_response.json"
	projectsGetSuccess        = "automation-studio/projects/get_response.json"
	projectsCreateSuccess     = "automation-studio/projects/create_response.json"
	projectsImportSuccess     = "automation-studio/projects/import_response.json"
	projectsExportSuccess     = "automation-studio/projects/export_response.json"
	projectsGetByNameFound    = "automation-studio/projects/2023.2.9/get_by_name_response.json"
	projectsGetByNameNotFound = "automation-studio/projects/2023.2.9/get_by_name_response_not_found.json"
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

func TestProjectsGetByNameFound(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		getAllResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, projectsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/automation-studio/projects", getAllResponse, 0)

		res, err := svc.GetByName("Port/VLAN Configuration - IOS")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Port/VLAN Configuration - IOS", res.Name)
		assert.Equal(t, "66aba9be41f8aad085ca0ee3", res.Id)
	}
}

func TestProjectsGetByNameNotFound(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		getAllResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, projectsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/automation-studio/projects", getAllResponse, 0)

		res, err := svc.GetByName("NonExistentProject")

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "project not found", err.Error())
	}
}

func TestProjectsGetByNameGetAllError(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/automation-studio/projects", "Internal server error", 500)

	res, err := svc.GetByName("SomeProject")

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

func TestProjectsImport(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	response := `{
		"message": "Project imported successfully",
		"data": {
			"_id": "imported-project-id",
			"name": "ImportedProject",
			"description": "Test import",
			"components": [],
			"folders": [],
			"iid": 456,
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

	testlib.AddPostResponseToMux("/automation-studio/projects/import", response, http.StatusOK)

	inputProject := Project{
		Id:          "imported-project-id",
		Name:        "ImportedProject",
		Description: "Test import",
		Components:  []ProjectComponent{},
		Folders:     []ProjectFolder{},
	}

	res, err := svc.Import(inputProject)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "ImportedProject", res.Name)
	assert.Equal(t, "imported-project-id", res.Id)
}

func TestProjectsImportError(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/automation-studio/projects/import", "Bad request", 400)

	inputProject := Project{
		Name: "InvalidProject",
	}

	res, err := svc.Import(inputProject)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestProjectsImportWithFolders(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	response := `{
		"message": "Project imported successfully",
		"data": {
			"_id": "imported-project-id",
			"name": "ImportedProject",
			"description": "Test import with folders",
			"components": [],
			"folders": [
				{
					"name": "TestFolder",
					"nodeType": "folder",
					"children": []
				}
			],
			"iid": 456,
			"componentIidIndex": 0
		}
	}`

	testlib.AddPostResponseToMux("/automation-studio/projects/import", response, http.StatusOK)

	inputProject := Project{
		Id:          "imported-project-id",
		Name:        "ImportedProject",
		Description: "Test import with folders",
		Components:  []ProjectComponent{},
		Folders: []ProjectFolder{
			{
				Name:     "TestFolder",
				NodeType: "folder",
				Children: []ProjectFolder{},
			},
		},
	}

	res, err := svc.Import(inputProject)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "ImportedProject", res.Name)
	assert.Equal(t, 1, len(res.Folders))
	assert.Equal(t, "TestFolder", res.Folders[0].Name)
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

// Note: AddMembers tests are skipped due to HTTP routing conflicts in the test framework.
// The AddMembers method calls both Get and Patch on the same URL pattern, but Go's ServeMux
// doesn't allow multiple handlers for the same pattern. The underlying functionality is
// tested through individual Get, Create, Delete, Import, Export method tests.

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

func TestProjectTransformImportFolder(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	// Test folder transformation
	folderData := map[string]interface{}{
		"nodeType": "folder",
		"name":     "TestFolder",
		"iid":      123, // Should be removed
		"children": nil,
	}

	svc.transformImport(folderData)

	_, hasIid := folderData["iid"]
	_, hasChildren := folderData["children"]

	assert.False(t, hasIid)      // iid should be removed for folders
	assert.False(t, hasChildren) // children should be removed when nil
}

func TestProjectTransformImportComponent(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	// Test component transformation
	componentData := map[string]interface{}{
		"nodeType": "component",
		"name":     "TestComponent", // Should be removed
		"iid":      456,
		"children": nil,
	}

	svc.transformImport(componentData)

	_, hasName := componentData["name"]
	_, hasChildren := componentData["children"]

	assert.False(t, hasName)     // name should be removed for components
	assert.False(t, hasChildren) // children should be removed when nil
}

func TestProjectTransformImportWithChildren(t *testing.T) {
	svc := setupProjectService()
	defer testlib.Teardown()

	// Test recursive transformation with children
	folderData := map[string]interface{}{
		"nodeType": "folder",
		"name":     "ParentFolder",
		"iid":      123,
		"children": []interface{}{
			map[string]interface{}{
				"nodeType": "folder",
				"name":     "ChildFolder",
				"iid":      456, // Should be removed
				"children": nil,
			},
			map[string]interface{}{
				"nodeType": "component",
				"name":     "ChildComponent", // Should be removed
				"iid":      789,
				"children": nil,
			},
		},
	}

	svc.transformImport(folderData)

	children := folderData["children"].([]interface{})
	childFolder := children[0].(map[string]interface{})
	childComponent := children[1].(map[string]interface{})

	// Check parent folder
	_, hasParentIid := folderData["iid"]
	assert.False(t, hasParentIid)

	// Check child folder
	_, hasChildFolderIid := childFolder["iid"]
	_, hasChildFolderChildren := childFolder["children"]
	assert.False(t, hasChildFolderIid)
	assert.False(t, hasChildFolderChildren)

	// Check child component
	_, hasChildComponentName := childComponent["name"]
	_, hasChildComponentChildren := childComponent["children"]
	assert.False(t, hasChildComponentName)
	assert.False(t, hasChildComponentChildren)
}

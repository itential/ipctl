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
	rolesGetSuccess    = "authorization/roles/get.success.json"
	rolesGetAllSuccess = "authorization/roles/getall.success.json"
	rolesCreateSuccess = "authorization/roles/create.success.json"
	rolesDeleteSuccess = "authorization/roles/delete.success.json"
	rolesImportSuccess = "authorization/roles/import.success.json"
)

func setupRoleService() *RoleService {
	return NewRoleService(
		testlib.Setup(),
	)
}

// TestNewRoleService tests the constructor
func TestNewRoleService(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewRoleService(client)

	assert.NotNil(t, svc)
	assert.Equal(t, client, svc.client)
}

// TestRoleServiceGetAll tests retrieving all roles
func TestRoleServiceGetAll(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, rolesGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/authorization/roles", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 2, len(res))
		assert.Equal(t, "admin_role", res[0].Name)
		assert.Equal(t, "user_role", res[1].Name)
		assert.Equal(t, "local_aaa", res[0].Provenance)
		assert.NotEmpty(t, res[0].Id)
		assert.NotEmpty(t, res[0].AllowedMethods)
		assert.NotEmpty(t, res[0].AllowedViews)
	}
}

// TestRoleServiceGetAllError tests error handling for GetAll
func TestRoleServiceGetAllError(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authorization/roles", "", 0)

	res, err := svc.GetAll()

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestRoleServiceGet tests retrieving a specific role by ID
func TestRoleServiceGet(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, rolesGetSuccess),
		)
		testlib.AddGetResponseToMux("/authorization/roles/{id}", response, 0)

		res, err := svc.Get("678ec1ab6d849669c6dbe952")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Role)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "admin_role", res.Name)
		assert.Equal(t, "local_aaa", res.Provenance)
		assert.Equal(t, "678ec1ab6d849669c6dbe952", res.Id)
		assert.Equal(t, "Administrative role with full permissions", res.Description)
		assert.Len(t, res.AllowedMethods, 3)
		assert.Len(t, res.AllowedViews, 2)
		assert.Equal(t, "createUser", res.AllowedMethods[0].Name)
		assert.Equal(t, "/admin/users", res.AllowedViews[0].Path)
	}
}

// TestRoleServiceGetError tests error handling for Get
func TestRoleServiceGetError(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	testlib.AddGetErrorToMux("/authorization/roles", "", 0)

	res, err := svc.Get("invalid-id")

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestRoleServiceCreate tests creating a new role
func TestRoleServiceCreate(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, rolesCreateSuccess),
		)
		testlib.AddPostResponseToMux("/authorization/roles", response, http.StatusOK)

		inputRole := Role{
			Name:        "test_role",
			Provenance:  "local_aaa",
			Description: "Test role for validation",
			AllowedMethods: []RoleMethod{
				{Name: "testMethod", Provenance: "local_aaa"},
			},
			AllowedViews: []RoleView{
				{Provenance: "local_aaa", Path: "/test"},
			},
		}

		res, err := svc.Create(inputRole)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test_role", res.Name)
		assert.Equal(t, "local_aaa", res.Provenance)
		assert.NotEmpty(t, res.Id)
		assert.Len(t, res.AllowedMethods, 1)
		assert.Len(t, res.AllowedViews, 1)
	}
}

// TestRoleServiceCreateWithEmptyArrays tests creating a role with empty allowed methods and views
func TestRoleServiceCreateWithEmptyArrays(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, rolesCreateSuccess),
		)
		testlib.AddPostResponseToMux("/authorization/roles", response, http.StatusOK)

		inputRole := Role{
			Name:           "test_role",
			Provenance:     "local_aaa",
			Description:    "Test role with empty arrays",
			AllowedMethods: []RoleMethod{},
			AllowedViews:   []RoleView{},
		}

		res, err := svc.Create(inputRole)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test_role", res.Name)
	}
}

// TestRoleServiceCreateError tests error handling for Create
func TestRoleServiceCreateError(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/authorization/roles", "", http.StatusInternalServerError)

	inputRole := Role{
		Name:       "test_role",
		Provenance: "local_aaa",
	}

	res, err := svc.Create(inputRole)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestRoleServiceDelete tests deleting a role
func TestRoleServiceDelete(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, rolesDeleteSuccess),
		)
		testlib.AddDeleteResponseToMux("/authorization/roles/{id}", response, http.StatusOK)

		err := svc.Delete("678ec1ab6d849669c6dbe952")

		assert.Nil(t, err)
	}
}

// TestRoleServiceDeleteError tests error handling for Delete
func TestRoleServiceDeleteError(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	testlib.AddDeleteErrorToMux("/authorization/roles/{id}", "", http.StatusNotFound)

	err := svc.Delete("invalid-id")

	assert.NotNil(t, err)
}

// TestRoleServiceImport tests importing a role
func TestRoleServiceImport(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, rolesImportSuccess),
		)
		testlib.AddPostResponseToMux("/authorization/roles", response, http.StatusOK)

		inputRole := Role{
			Id:          "678ec1ab6d849669c6dbe955",
			Name:        "imported_role",
			Provenance:  "local_aaa",
			Description: "Role imported from external source",
			AllowedMethods: []RoleMethod{
				{Name: "importMethod", Provenance: "local_aaa"},
			},
			AllowedViews: []RoleView{
				{Provenance: "local_aaa", Path: "/import"},
			},
		}

		res, err := svc.Import(inputRole)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "imported_role", res.Name)
		assert.Equal(t, "local_aaa", res.Provenance)
		assert.Equal(t, "678ec1ab6d849669c6dbe955", res.Id)
		assert.Len(t, res.AllowedMethods, 1)
		assert.Len(t, res.AllowedViews, 1)
	}
}

// TestRoleServiceImportError tests error handling for Import
func TestRoleServiceImportError(t *testing.T) {
	svc := setupRoleService()
	defer testlib.Teardown()

	testlib.AddPostErrorToMux("/authorization/roles", "", http.StatusInternalServerError)

	inputRole := Role{
		Name:       "imported_role",
		Provenance: "local_aaa",
	}

	res, err := svc.Import(inputRole)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

// TestRoleStructFields tests the Role struct field mappings
func TestRoleStructFields(t *testing.T) {
	role := Role{
		Id:          "test-id",
		Name:        "test-role",
		Provenance:  "local_aaa",
		Description: "Test description",
		AllowedMethods: []RoleMethod{
			{Name: "testMethod", Provenance: "local_aaa"},
		},
		AllowedViews: []RoleView{
			{Provenance: "local_aaa", Path: "/test"},
		},
	}

	assert.Equal(t, "test-id", role.Id)
	assert.Equal(t, "test-role", role.Name)
	assert.Equal(t, "local_aaa", role.Provenance)
	assert.Equal(t, "Test description", role.Description)
	assert.Len(t, role.AllowedMethods, 1)
	assert.Len(t, role.AllowedViews, 1)
	assert.Equal(t, "testMethod", role.AllowedMethods[0].Name)
	assert.Equal(t, "/test", role.AllowedViews[0].Path)
}

// TestRoleMethodStruct tests the RoleMethod struct
func TestRoleMethodStruct(t *testing.T) {
	method := RoleMethod{
		Name:       "testMethod",
		Provenance: "local_aaa",
	}

	assert.Equal(t, "testMethod", method.Name)
	assert.Equal(t, "local_aaa", method.Provenance)
}

// TestRoleViewStruct tests the RoleView struct
func TestRoleViewStruct(t *testing.T) {
	view := RoleView{
		Provenance: "local_aaa",
		Path:       "/admin/test",
	}

	assert.Equal(t, "local_aaa", view.Provenance)
	assert.Equal(t, "/admin/test", view.Path)
}

// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	prebuiltsGetAllSuccess = "prebuilts/getall.success.json"
	prebuiltsGetSuccess    = "prebuilts/get.success.json"
	prebuiltsGetNotFound   = "prebuilts/get.notfound.json"
	prebuiltsImportSuccess = "prebuilts/import.success.json"
	prebuiltsExportSuccess = "prebuilts/export.success.json"
)

func setupPrebuiltService() *PrebuiltService {
	return NewPrebuiltService(
		testlib.Setup(),
	)
}

func TestPrebuiltsGetAll(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/prebuilts", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))

		// Verify first prebuilt
		assert.Equal(t, "prebuilt-1", res[0].Id)
		assert.Equal(t, "Test Prebuilt 1", res[0].Name)
		assert.Equal(t, "Test Author", res[0].Author)
		assert.Equal(t, "1.0.0", res[0].Version)
		assert.Equal(t, "MIT", res[0].License)
		assert.Contains(t, res[0].Keywords, "test")
		assert.Contains(t, res[0].Keywords, "automation")
		assert.Equal(t, 1, len(res[0].Components))

		// Verify second prebuilt
		assert.Equal(t, "prebuilt-2", res[1].Id)
		assert.Equal(t, "Test Prebuilt 2", res[1].Name)
		assert.Equal(t, "Another Author", res[1].Author)
		assert.Equal(t, "2.1.0", res[1].Version)
		assert.Equal(t, "Apache-2.0", res[1].License)
		assert.Equal(t, 2, len(res[1].Components))
	}
}

func TestPrebuiltsGet(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsGetSuccess),
		)

		testlib.AddGetResponseToMux("/prebuilts/prebuilt-1", response, 0)

		res, err := svc.Get("prebuilt-1")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "prebuilt-1", res.Id)
		assert.Equal(t, "Test Prebuilt 1", res.Name)
		assert.Equal(t, "Test Author", res.Author)
		assert.Equal(t, "1.0.0", res.Version)
		assert.Equal(t, "MIT", res.License)
		assert.Equal(t, "Test prebuilt for unit testing", res.Description)

		// Verify manifest
		assert.Equal(t, "test-prebuilt-1", res.Manifest.BundleName)
		assert.Equal(t, "abc123def456", res.Manifest.Fingerprint)
		assert.Equal(t, 1, len(res.Manifest.Artifacts))

		// Verify repository
		assert.Equal(t, "github.com", res.Repository.Hostname)
		assert.Equal(t, "/test/prebuilt-1", res.Repository.Path)
		assert.Equal(t, "git", res.Repository.Type)

		// Verify components
		assert.Equal(t, 1, len(res.Components))
		assert.Equal(t, "component-1", res.Components[0].Id)
		assert.Equal(t, "Test Component 1", res.Components[0].Name)
		assert.Equal(t, "workflow", res.Components[0].Type)

		// Verify dependencies
		assert.NotNil(t, res.Dependencies)
		assert.Contains(t, res.Dependencies, "@itential/core")
	}
}

func TestPrebuiltsGetNotFound(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsGetNotFound),
		)

		testlib.AddGetResponseToMux("/prebuilts/nonexistent-id", response, http.StatusNotFound)

		res, err := svc.Get("nonexistent-id")

		assert.NotNil(t, err)
		assert.Nil(t, res)
	}
}

func TestPrebuiltsGetByName(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		// First mock the GetAll call
		getAllResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsGetAllSuccess),
		)
		testlib.AddGetResponseToMux("/prebuilts", getAllResponse, 0)

		// Then mock the Get call for the specific ID
		getResponse := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsGetSuccess),
		)
		testlib.AddGetResponseToMux("/prebuilts/prebuilt-1", getResponse, 0)

		res, err := svc.GetByName("Test Prebuilt 1")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "prebuilt-1", res.Id)
		assert.Equal(t, "Test Prebuilt 1", res.Name)
	}
}

func TestPrebuiltsGetByNameNotFound(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsGetAllSuccess),
		)
		testlib.AddGetResponseToMux("/prebuilts", response, 0)

		res, err := svc.GetByName("Nonexistent Prebuilt")

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "prebuilt `Nonexistent Prebuilt` does not exist")
	}
}

func TestPrebuiltsDelete(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	for range fixtureSuites {
		testlib.AddDeleteResponseToMux("/prebuilts/prebuilt-1", "", 0)

		err := svc.Delete("prebuilt-1")

		assert.Nil(t, err)
	}
}

func TestPrebuiltsImport(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	// Create test prebuilt package
	testPackage := PrebuiltPackage{
		Metadata: PrebuiltPackageMetadata{
			Name:         "Test Import Package",
			Version:      "1.0.0",
			Description:  "Package for import testing",
			License:      "MIT",
			Author:       "Test Author",
			Keywords:     []string{"test", "import"},
			Dependencies: map[string]interface{}{"@itential/core": "^2023.2.0"},
			GitlabId:     12345,
			Repository: PrebuiltPackageMetadataRepository{
				Type:     "git",
				Hostname: "github.com",
				Path:     "/test/import-package",
			},
		},
		Manifest: PrebuiltPackageManifest{
			Name:        "test-import-package",
			Fingerprint: "import123abc",
			Epoch:       "1705488900",
			Artifacts: []PrebuiltPackageManifestArtifact{
				{
					Id:       "import-artifact-1",
					Name:     "Import Test Workflow",
					Type:     "workflow",
					Location: "/workflows/import-test",
					Properties: map[string]interface{}{
						"type": "workflow",
					},
				},
			},
		},
		Bundles: []PrebuiltPackageBundle{
			{
				Type: "workflows",
				Data: map[string]interface{}{
					"workflow-1": map[string]interface{}{
						"name": "Import Test Workflow",
						"type": "workflow",
					},
				},
			},
		},
		Readme: "# Import Test Package\n\nThis is a test package.",
	}

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsImportSuccess),
		)

		testlib.AddPostResponseToMux("/prebuilts/import", response, http.StatusOK)

		res, err := svc.Import(testPackage, false)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "prebuilt-imported", res.Id)
		assert.Equal(t, "Imported Prebuilt", res.Name)
		assert.Equal(t, "Import Test Author", res.Author)
	}
}

func TestPrebuiltsImportWithOverwrite(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	testPackage := PrebuiltPackage{
		Metadata: PrebuiltPackageMetadata{
			Name:    "Test Overwrite Package",
			Version: "2.0.0",
		},
	}

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsImportSuccess),
		)

		testlib.AddPostResponseToMux("/prebuilts/import", response, http.StatusOK)

		res, err := svc.Import(testPackage, true)

		assert.Nil(t, err)
		assert.NotNil(t, res)
	}
}

func TestPrebuiltsImportRaw(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	rawData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":    "Raw Import Test",
			"version": "1.0.0",
		},
		"manifest": map[string]interface{}{
			"bundleName": "raw-import-test",
		},
	}

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsImportSuccess),
		)

		testlib.AddPostResponseToMux("/prebuilts/import", response, http.StatusOK)

		res, err := svc.ImportRaw(rawData, false)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "prebuilt-imported", res.Id)
	}
}

func TestPrebuiltsExport(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, prebuiltsExportSuccess),
		)

		testlib.AddGetResponseToMux("/prebuilts/prebuilt-1/export", response, 0)

		res, err := svc.Export("prebuilt-1")

		assert.Nil(t, err)
		assert.NotNil(t, res)

		// Verify metadata
		assert.Equal(t, "Test Prebuilt Package", res.Metadata.Name)
		assert.Equal(t, "1.0.0", res.Metadata.Version)
		assert.Equal(t, "Test prebuilt package for export testing", res.Metadata.Description)
		assert.Equal(t, "MIT", res.Metadata.License)
		assert.Equal(t, "Export Test Author", res.Metadata.Author)
		assert.Equal(t, 12345, res.Metadata.GitlabId)

		// Verify repository
		assert.Equal(t, "git", res.Metadata.Repository.Type)
		assert.Equal(t, "github.com", res.Metadata.Repository.Hostname)
		assert.Equal(t, "/test/prebuilt-package", res.Metadata.Repository.Path)

		// Verify manifest
		assert.Equal(t, "test-prebuilt-package", res.Manifest.Name)
		assert.Equal(t, "export123def456", res.Manifest.Fingerprint)
		assert.Equal(t, 1, len(res.Manifest.Artifacts))

		// Verify artifacts
		artifact := res.Manifest.Artifacts[0]
		assert.Equal(t, "export-artifact-1", artifact.Id)
		assert.Equal(t, "Export Test Workflow", artifact.Name)
		assert.Equal(t, "workflow", artifact.Type)
		assert.Equal(t, "/workflows/export-test-workflow", artifact.Location)

		// Verify bundles
		assert.Equal(t, 1, len(res.Bundles))
		assert.Equal(t, "workflows", res.Bundles[0].Type)
		assert.NotNil(t, res.Bundles[0].Data)

		// Verify readme
		assert.Contains(t, res.Readme, "# Test Prebuilt Package")
		assert.Contains(t, res.Readme, "export functionality testing")
	}
}

func TestNewPrebuiltService(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewPrebuiltService(client)

	assert.NotNil(t, svc)
	assert.Equal(t, client, svc.client)
}

func TestPrebuiltServiceMethodsExist(t *testing.T) {
	svc := setupPrebuiltService()
	defer testlib.Teardown()

	// Verify all expected methods exist and have correct signatures
	assert.NotNil(t, svc.Get)
	assert.NotNil(t, svc.GetByName)
	assert.NotNil(t, svc.GetAll)
	assert.NotNil(t, svc.Delete)
	assert.NotNil(t, svc.Import)
	assert.NotNil(t, svc.ImportRaw)
	assert.NotNil(t, svc.Export)
}

func TestPrebuiltTypeValidation(t *testing.T) {
	// Test that all struct types can be properly marshaled/unmarshaled
	testCases := []struct {
		name string
		data interface{}
	}{
		{
			name: "Prebuilt",
			data: Prebuilt{
				Id:      "test-id",
				Name:    "Test Name",
				Author:  "Test Author",
				Version: "1.0.0",
			},
		},
		{
			name: "PrebuiltPackage",
			data: PrebuiltPackage{
				Metadata: PrebuiltPackageMetadata{
					Name:    "Test Package",
					Version: "1.0.0",
				},
			},
		},
		{
			name: "PrebuiltComponent",
			data: PrebuiltComponent{
				Id:   "component-id",
				Name: "Component Name",
				Type: "workflow",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tc.data)
			require.NoError(t, err)
			assert.NotEmpty(t, jsonData)

			// Test JSON unmarshaling
			switch tc.data.(type) {
			case Prebuilt:
				var unmarshaled Prebuilt
				err = json.Unmarshal(jsonData, &unmarshaled)
				require.NoError(t, err)
			case PrebuiltPackage:
				var unmarshaled PrebuiltPackage
				err = json.Unmarshal(jsonData, &unmarshaled)
				require.NoError(t, err)
			case PrebuiltComponent:
				var unmarshaled PrebuiltComponent
				err = json.Unmarshal(jsonData, &unmarshaled)
				require.NoError(t, err)
			}
		})
	}
}

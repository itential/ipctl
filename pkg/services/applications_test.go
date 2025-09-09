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
	applicationsGetAllSuccess  = "applications/getall.success.json"
	applicationsGetSuccess     = "applications/get.success.json"
	applicationsCreateSuccess  = "applications/create.success.json"
	applicationsStartSuccess   = "applications/start.success.json"
	applicationsStopSuccess    = "applications/stop.success.json"
	applicationsRestartSuccess = "applications/restart.success.json"
)

func setupApplicationService() *ApplicationService {
	return NewApplicationService(
		testlib.Setup(),
	)
}

func TestApplicationsGetAll(t *testing.T) {
	svc := setupApplicationService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, applicationsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/applications", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 18, len(res))
	}
}

func TestApplicationsGet(t *testing.T) {
	svc := setupApplicationService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, applicationsGetSuccess),
		)

		testlib.AddGetResponseToMux("/applications/WorkFlowEngine", response, 0)

		res, err := svc.Get("WorkFlowEngine")

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Application)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "WorkFlowEngine", res.Name)
	}
}

func TestNewApplicationService(t *testing.T) {
	client := testlib.Setup()
	defer testlib.Teardown()

	svc := NewApplicationService(client)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.client)
	assert.Equal(t, reflect.TypeOf((*ApplicationService)(nil)), reflect.TypeOf(svc))
}

func TestApplicationsCreate(t *testing.T) {
	svc := setupApplicationService()
	defer testlib.Teardown()

	app := Application{
		Name:  "TestApplication",
		Type:  "service",
		Model: "App",
		Properties: map[string]interface{}{
			"version":     "1.0.0",
			"description": "Test application",
		},
		IsEncrypted: false,
		LoggerProperties: map[string]interface{}{
			"level": "info",
		},
	}

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, applicationsCreateSuccess),
		)

		testlib.AddPostResponseToMux("/applications", response, http.StatusOK)

		res, err := svc.Create(app)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reflect.TypeOf((*Application)(nil)), reflect.TypeOf(res))
		assert.Equal(t, "TestApplication", res.Name)
	}
}

func TestApplicationsStart(t *testing.T) {
	svc := setupApplicationService()
	defer testlib.Teardown()

	appName := "TestApplication"

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, applicationsStartSuccess),
		)

		testlib.AddPutResponseToMux("/applications/TestApplication/start", response, 0)

		err := svc.Start(appName)

		assert.Nil(t, err)
	}
}

func TestApplicationsStop(t *testing.T) {
	svc := setupApplicationService()
	defer testlib.Teardown()

	appName := "TestApplication"

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, applicationsStopSuccess),
		)

		testlib.AddPutResponseToMux("/applications/TestApplication/stop", response, 0)

		err := svc.Stop(appName)

		assert.Nil(t, err)
	}
}

func TestApplicationsRestart(t *testing.T) {
	svc := setupApplicationService()
	defer testlib.Teardown()

	appName := "TestApplication"

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, applicationsRestartSuccess),
		)

		testlib.AddPutResponseToMux("/applications/TestApplication/restart", response, 0)

		err := svc.Restart(appName)

		assert.Nil(t, err)
	}
}

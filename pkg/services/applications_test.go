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
	applicationsGetAllSuccess = "applications/getall.success.json"
	applicationsGetSuccess    = "applications/get.success.json"
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

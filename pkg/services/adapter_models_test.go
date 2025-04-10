// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"path/filepath"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	adapterModelsGetAllSuccess = "adapter-models/getall.success.json"
)

func setupAdapterModelService() *AdapterModelService {
	return NewAdapterModelService(
		testlib.Setup(),
	)
}

func TestAdapterModelsGetAll(t *testing.T) {
	svc := setupAdapterModelService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, adapterModelsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/adapter-models/types", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, 7, len(res))
	}
}

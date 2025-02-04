// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	adapterModelsGetAllResponse = testlib.Fixture("testdata/adapter-models/getall.json")
)

func setupAdapterModelService() *AdapterModelService {
	return NewAdapterModelService(
		testlib.Setup(),
	)
}

func TestAdapterModelsGetAll(t *testing.T) {
	svc := setupAdapterModelService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/adapter-models/types", adapterModelsGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 7, len(res))
}

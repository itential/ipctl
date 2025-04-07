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
	methodsGetAllSuccess = "authorization/methods/getall.success.json"
)

func setupMethodService() *MethodService {
	return NewMethodService(
		testlib.Setup(),
	)
}

func TestMethodsGetAll(t *testing.T) {
	svc := setupMethodService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, methodsGetAllSuccess),
		)

		testlib.AddGetResponseToMux("/authorization/methods", response, 0)

		res, err := svc.GetAll()

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 523, len(res))
	}
}

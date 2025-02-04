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
	methodsGetAllResponse = testlib.Fixture("testdata/authorization/methods/getall.json")
)

func setupMethodService() *MethodService {
	return NewMethodService(
		testlib.Setup(),
	)
}

func TestMethodsGetAll(t *testing.T) {
	svc := setupMethodService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/authorization/methods", methodsGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 523, len(res))
}

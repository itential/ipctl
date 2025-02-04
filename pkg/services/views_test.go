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
	viewsGetAllResponse = testlib.Fixture("testdata/authorization/views/getall.json")
)

func setupViewService() *ViewService {
	return NewViewService(
		testlib.Setup(),
	)
}

func TestViewsGetAll(t *testing.T) {
	svc := setupViewService()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/authorization/views", viewsGetAllResponse, 0)

	res, err := svc.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 51, len(res))
}

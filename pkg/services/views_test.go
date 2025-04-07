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
	viewsGetAllSuccess = "authorization/views/getall.success.json"
)

func setupViewService() *ViewService {
	return NewViewService(
		testlib.Setup(),
	)
}

func TestViewsGetAll(t *testing.T) {
	svc := setupViewService()
	defer testlib.Teardown()

	for _, ele := range fixtureSuites {
		response := testlib.Fixture(
			filepath.Join(fixtureRoot, ele, viewsGetAllSuccess),
		)

		data, err := fixtureDataToMap(response)
		if err != nil {
			t.FailNow()
		}

		testlib.AddGetResponseToMux("/authorization/views", response, 0)

		res, err := svc.GetAll()

		total := data["total"].(float64)

		assert.Nil(t, err)
		assert.Equal(t, int(total), len(res))
	}
}

// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	accountsGetResponse    = testlib.Fixture("testdata/authorization/accounts/get.json")
	accountsGetAllResponse = testlib.Fixture("testdata/authorization/accounts/getall.json")
)

func TestGet(t *testing.T) {
	runner := NewAccountRunner(
		testlib.Setup(),
		testlib.DefaultConfig(),
	)

	testlib.AddGetResponseToMux("/authorization/accounts", accountsGetAllResponse, 0)

	res, err := runner.Get(Request{})

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Empty(t, res.Text)
	assert.NotEmpty(t, res.Lines)
	assert.NotEmpty(t, res.Json)
}

func TestDescribe(t *testing.T) {
	runner := NewAccountRunner(
		testlib.Setup(),
		testlib.DefaultConfig(),
	)

	testlib.AddGetResponseToMux("/authorization/accounts", accountsGetAllResponse, 0)

	res, err := runner.Get(Request{})

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Empty(t, res.Text)
	assert.NotEmpty(t, res.Lines)
	assert.NotEmpty(t, res.Json)
}

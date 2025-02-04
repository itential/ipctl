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
	getCurrentUserResponse = testlib.Fixture("testdata/whoami.json")
)

func TestGetCurrentUser(t *testing.T) {
	testClient := testlib.Setup()
	defer testlib.Teardown()

	testlib.AddGetResponseToMux("/whoami", getCurrentUserResponse, 0)

	res, err := GetCurrentUser(testClient)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "admin@pronghorn", res.Username)
}

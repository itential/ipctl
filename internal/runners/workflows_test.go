// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"net/http"
	"testing"

	"github.com/itential/ipctl/internal/testlib"
	"github.com/stretchr/testify/assert"
)

var (
	workflowsDeleteResponse = testlib.Fixture("testdata/workflows/delete.success.json")
)

func TestWorkflowDelete(t *testing.T) {
	runner := NewWorkflowRunner(
		testlib.Setup(),
		testlib.DefaultConfig(),
	)

	testlib.AddDeleteResponseToMux(
		"/workflow_builder/workflows/delete/{name}",
		workflowsDeleteResponse,
		http.StatusOK,
	)

	res, err := runner.Delete(Request{
		Args: []string{"cli-test-1"},
	})

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Text)
}

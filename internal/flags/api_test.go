// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"testing"
)

func TestApiPutOptions(t *testing.T) {
	checkFlags(t, &ApiPutOptions{}, []string{"data", "expected-status-code"})
}

func TestApiDeleteOptions(t *testing.T) {
	checkFlags(t, &ApiDeleteOptions{}, []string{"expected-status-code"})
}

func TestApiPostOptions(t *testing.T) {
	checkFlags(t, &ApiPostOptions{}, []string{"data", "expected-status-code"})
}

func TestApiPatchOptions(t *testing.T) {
	checkFlags(t, &ApiPatchOptions{}, []string{"data", "expected-status-code"})
}

// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "testing"

func TestRoleCreateOptions(t *testing.T) {
	checkFlags(t, &RoleCreateOptions{}, []string{"method", "view"})
}

func TestRoleGetOptions(t *testing.T) {
	checkFlags(t, &RoleGetOptions{}, []string{"all"})
}

func TestRoleDescribeOptions(t *testing.T) {
	checkFlags(t, &RoleDescribeOptions{}, []string{"type"})
}

func TestRoleExportOptions(t *testing.T) {
	checkFlags(t, &RoleExportOptions{}, []string{"type"})
}

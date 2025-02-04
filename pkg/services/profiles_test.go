// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import "github.com/itential/ipctl/internal/testlib"

var (
	profileCreateExistsResponse   = testlib.Fixture("testdata/profiles/create.exists.json")
	profileCreateResponse         = testlib.Fixture("testdata/profiles/create.json")
	profileDeleteResponse         = testlib.Fixture("testdata/profiles/delete.json")
	profileDeleteNotFoundResponse = testlib.Fixture("testdata/profiles/delete.notfound.json")
	profileExportResponse         = testlib.Fixture("testdata/profiles/export.json")
	profileExportNotFoundResponse = testlib.Fixture("testdata/profiles/export.notfound.json")
	profileGetAllResponse         = testlib.Fixture("testdata/profiles/getall.json")
	profileGetResponse            = testlib.Fixture("testdata/profiles/get.json")
	profileGetNotFoundResponse    = testlib.Fixture("testdata/profiles/get.notfound.json")
	profileImportExistsResponse   = testlib.Fixture("testdata/profiles/import.exists.json")
	profileImportResponse         = testlib.Fixture("testdata/profiles/import.json")
)

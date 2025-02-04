// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"encoding/json"
)

// Unmarshal will take in a mapped interface object and load it into the a
// instance of a struct specified by ptr
func Unmarshal(in any, ptr any) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, ptr)
}

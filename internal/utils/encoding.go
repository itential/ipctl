// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package utils

import (
	"encoding/json"

	"github.com/itential/ipctl/pkg/logger"
	"gopkg.in/yaml.v2"
)

// ToMap accepts any object and will return at as a map using json marshal and
// unmarshal.  This fuction will return an error if if fails to marshal or
// unmarshal the input object.
func ToMap(in any, out any) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &out)
}

// UnmarshalData will attempt to to unmarshal a byte array into an object.  It
// will first attempt to unmarshal the byte array as JSON.  If that fails, it
// will attempt to unmarshal the data as YAML.
func UnmarshalData(data []byte, ptr any) {
	// FIXME (privateip) This function should be refactored to return an error
	// instead of simply logging a fatal error.
	if err := json.Unmarshal(data, ptr); err != nil {
		logger.Error(err, "failed to unmarshal json data")
		if err = yaml.Unmarshal(data, ptr); err != nil {
			logger.Fatal(err, "failed to unmarshal data")
		}
	}
}

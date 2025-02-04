// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package utils

import (
	"encoding/json"

	"github.com/itential/ipctl/pkg/logger"
	"gopkg.in/yaml.v2"
)

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
	if err := json.Unmarshal(data, ptr); err != nil {
		logger.Error(err, "failed to unmarshal json data")
		if err = yaml.Unmarshal(data, ptr); err != nil {
			logger.Fatal(err, "failed to unmarshal data")
		}
	}
}

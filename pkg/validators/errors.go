// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package validators

const (
	ValidatorUnknownError = 100 // Default error
	ValidatorServiceError = 101 // Error returned from the service
	ValidatorAssetExists  = 102 // Asset exists on the server
)

type ValidatorError struct {
	// The validation code returned as part of the validation error message
	// which gives an indication to the calling function as to the general
	// nature of the validation error.
	Code int

	// The error returned from the function.
	Err error
}

func (v ValidatorError) Error() string {
	return v.Err.Error()
}

func NewValidatorError(code int, e error) *ValidatorError {
	if code == 0 {
		code = ValidatorUnknownError
	}
	return &ValidatorError{
		Code: code,
		Err:  e,
	}
}

func ServiceError(e error) *ValidatorError {
	return NewValidatorError(ValidatorServiceError, e)
}

func AssetExists(e error) *ValidatorError {
	return NewValidatorError(ValidatorAssetExists, e)
}

package action

import "errors"

var (
	ErrConfigFileNotFound = errors.New("provided config file not found")
	ErrConfigFlagMissing  = errors.New("missing required env MEKONG_CONFIG_FILE")

	// Validation errors
	ErrValidateFieldNotExistsBackendHost = errors.New("key backendHost is missing in the provided config file")

	// Parsing error
	ErrHTTPMethodNotValid = errors.New("http method is not valid")

	// Rules error
	ErrRuleHasNoBody = errors.New("request needs to have a body")
	ErrRuleHasBody   = errors.New("request needs to have an empty body")

	ErrRuleHasNoQueryStrings = errors.New("request needs to have query strings")
	ErrRuleHasQueryStrings   = errors.New("request must not have query strings")
)

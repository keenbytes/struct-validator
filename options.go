package validator

// ValidateOptions is an optional configuration for validation:
// * RestrictFields defines what struct fields should be validated
// * TagName sets tag used to define validation (default is "validation")
// * OverwriteValues allows overriding values of fields
type ValidateOptions struct {
	RestrictFields  map[string]bool
	TagName         string
	OverwriteValues map[string]interface{}
}

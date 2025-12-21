package validator

// Optional configuration for validation:
// * RestrictFields defines what struct fields should be validated
// * TagName sets tag used to define validation (default is "validation")
type ValidateOptions struct {
	RestrictFields map[string]bool
	TagName        string
}

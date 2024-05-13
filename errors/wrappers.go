package errors

// Wrap std errors package
// we do this should we only need to import a signle errors package (making name conflicst easier)

import "errors"

// These variables are used to give us access to existing
// functions in the std lib errors package. We can also
// wrap them in custom functionality as needed if we want,
// or mock them during testing
var (
	As = errors.As
	Is = errors.Is
)

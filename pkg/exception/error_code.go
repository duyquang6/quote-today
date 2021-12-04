// Package exception provides error code definition
package exception

// Common module (00) error codes definition.
var (
	ErrValidation     = 4000001
	ErrUnauthorized   = 4010001
	ErrInternalServer = 5000001
	ErrUnknownError   = 500
)

// DateQuote module (01) error codes definition.
var (
	ErrNoQuoteFound = 4040101
)

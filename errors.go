package shopadmin

import "errors"

var (
	// ErrStoreRequired is returned when Store is not provided
	ErrStoreRequired = errors.New("store is required")

	// ErrLoggerRequired is returned when Logger is not provided
	ErrLoggerRequired = errors.New("logger is required")

	// ErrRegistryRequired is returned when Registry is not provided
	// (kept for backward compatibility with Routes())
	ErrRegistryRequired = errors.New("registry is required")
)

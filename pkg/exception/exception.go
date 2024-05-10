package exception

// Code is a type alias for string, representing the error code of an exception.
type Code string

// Predefined error codes.
const (
	InvalidArgumentCode  Code = "INVALID_ARGUMENT"  // Represents an invalid argument error.
	NotFoundCode         Code = "NOT_FOUND"         // Represents a not found error.
	AlreadyExistsCode    Code = "ALREADY_EXISTS"    // Represents an already exists error.
	PermissionDeniedCode Code = "PERMISSION_DENIED" // Represents a permission denied error.
	UnauthenticatedCode  Code = "UNAUTHENTICATED"   // Represents an unauthenticated error.
	InternalErrorCode    Code = "INTERNAL"          // Represents an internal error.
)

// Exception is a struct to represent exception/error from service.
// Code is the error code of the exception.
// Message is the error message of the exception.
// Error is the original error that caused the exception, if any.
type Exception struct {
	Code    Code
	Message any
	Error   error
}

// InvalidArgument creates a new Exception with the InvalidArgumentCode error code.
func InvalidArgument(message any) *Exception {
	return &Exception{
		Code:    InvalidArgumentCode,
		Message: message,
	}
}

// NotFound creates a new Exception with the NotFoundCode error code.
func NotFound(message any) *Exception {
	return &Exception{
		Code:    NotFoundCode,
		Message: message,
	}
}

// AlreadyExists creates a new Exception with the AlreadyExistsCode error code.
func AlreadyExists(message any) *Exception {
	return &Exception{
		Code:    AlreadyExistsCode,
		Message: message,
	}
}

// PermissionDenied creates a new Exception with the PermissionDeniedCode error code.
func PermissionDenied(message any) *Exception {
	return &Exception{
		Code:    PermissionDeniedCode,
		Message: message,
	}
}

// Unauthenticated creates a new Exception with the UnauthenticatedCode error code.
func Unauthenticated(message any) *Exception {
	return &Exception{
		Code:    UnauthenticatedCode,
		Message: message,
	}
}

// Internal creates a new Exception with the InternalErrorCode error code.
// The original error that caused the exception is also included.
func Internal(message any, err error) *Exception {
	return &Exception{
		Code:    InternalErrorCode,
		Message: message,
		Error:   err,
	}
}

// Conflict creates a new Exception with the AlreadyExistsCode error code.
func Conflict(message any) *Exception {
	return &Exception{
		Code:    AlreadyExistsCode,
		Message: message,
	}
}

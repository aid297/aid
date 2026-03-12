package driver

import (
	"errors"
	"fmt"

	"github.com/aid297/aid/simpleDB/kernal"
)

type ErrorCode string

const (
	ErrorCodeInvalidArgument ErrorCode = "invalid_argument"
	ErrorCodeNotFound        ErrorCode = "not_found"
	ErrorCodeConflict        ErrorCode = "conflict"
	ErrorCodeClosed          ErrorCode = "closed"
	ErrorCodeReadOnly        ErrorCode = "read_only"
	ErrorCodeUnauthorized    ErrorCode = "unauthorized"
	ErrorCodeDDL             ErrorCode = "ddl_error"
	ErrorCodeInternal        ErrorCode = "internal"
)

type DriverError struct {
	Code ErrorCode
	Err  error
}

func (e *DriverError) Error() string {
	if e == nil || e.Err == nil {
		return "driver error"
	}
	return fmt.Sprintf("driver[%s]: %v", e.Code, e.Err)
}

func (e *DriverError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*DriverError); ok {
		return err
	}

	switch {
	case errors.Is(err, kernal.ErrEmptyKey),
		errors.Is(err, kernal.ErrInvalidSchema),
		errors.Is(err, kernal.ErrInvalidQueryCondition),
		errors.Is(err, kernal.ErrBatchEmpty),
		errors.Is(err, kernal.ErrInvalidRegistration),
		errors.Is(err, kernal.ErrInvalidPermissionAssignment),
		errors.Is(err, kernal.ErrInvalidInitPassword),
		errors.Is(err, kernal.ErrInvalidTableAccessGrant):
		return &DriverError{Code: ErrorCodeInvalidArgument, Err: err}
	case errors.Is(err, kernal.ErrKeyNotFound),
		errors.Is(err, kernal.ErrRelationNotFound),
		errors.Is(err, kernal.ErrUserNotFound),
		errors.Is(err, kernal.ErrRoleNotFound),
		errors.Is(err, kernal.ErrTableOwnerNotFound):
		return &DriverError{Code: ErrorCodeNotFound, Err: err}
	case errors.Is(err, kernal.ErrInvalidCredentials),
		errors.Is(err, kernal.ErrUserInactive),
		errors.Is(err, kernal.ErrTableAccessDenied):
		return &DriverError{Code: ErrorCodeUnauthorized, Err: err}
	case errors.Is(err, kernal.ErrUniqueConflict),
		errors.Is(err, kernal.ErrPrimaryKeyConflict),
		errors.Is(err, kernal.ErrTxConflict),
		errors.Is(err, kernal.ErrUserAlreadyExists),
		errors.Is(err, kernal.ErrTableOwnerAlreadyAssigned):
		return &DriverError{Code: ErrorCodeConflict, Err: err}
	case errors.Is(err, kernal.ErrSuperAdminRoleReserved):
		return &DriverError{Code: ErrorCodeUnauthorized, Err: err}
	case errors.Is(err, kernal.ErrDatabaseClosed),
		errors.Is(err, kernal.ErrTxClosed):
		return &DriverError{Code: ErrorCodeClosed, Err: err}
	case errors.Is(err, kernal.ErrTxReadOnly):
		return &DriverError{Code: ErrorCodeReadOnly, Err: err}
	case errors.Is(err, kernal.ErrAlterTableInvalid),
		errors.Is(err, kernal.ErrColumnAlreadyExists),
		errors.Is(err, kernal.ErrColumnNotFound),
		errors.Is(err, kernal.ErrCannotDropPrimaryKey),
		errors.Is(err, kernal.ErrSystemBootstrap),
		errors.Is(err, kernal.ErrSystemTableSchema),
		errors.Is(err, kernal.ErrSchemaAlreadyExists),
		errors.Is(err, kernal.ErrSchemaNotConfigured):
		return &DriverError{Code: ErrorCodeDDL, Err: err}
	default:
		return &DriverError{Code: ErrorCodeInternal, Err: err}
	}
}

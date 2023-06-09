package validation

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun/driver/pgdriver"
	"strings"
)

var (
	ErrConstraintViolation     = fmt.Errorf("record does not satisfy some of the column constraints")
	ErrUniqConstraintViolation = fmt.Errorf("%w: some unique columns have duplicates", ErrConstraintViolation)
	ErrTimeout                 = fmt.Errorf("connection timed out")
	ErrInvalidEntity           = fmt.Errorf("field is not valid")
	ErrNil                     = fmt.Errorf("%w: field cannot be nil", ErrInvalidEntity)
	ErrNotAllowed              = fmt.Errorf("%w: the value is not allowed", ErrInvalidEntity)
	ErrMissingRelation         = fmt.Errorf("%w: a required relation is missing", ErrInvalidEntity)
	ErrInvalidCredentials      = fmt.Errorf("the given credentials are not valid")
	ErrValidated               = fmt.Errorf("the current link has already been validated")
	ErrNotFound                = fmt.Errorf("could not find any record matching the request")
	ErrUnauthorized            = fmt.Errorf("you are not allowed to perform this action")
)

// HandlePGError extends pg library typed errors. Only a few errors are typed to be targeted with errors.Is, and some
// pretty common errors aren't. This handler parses postgres errors in a more test-friendly way.
func HandlePGError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	pgErr, ok := err.(pgdriver.Error)
	if ok {
		if pgErr.IntegrityViolation() {
			// This error has a special treatment because, in most case, it is not checked upfront by the service.
			// Other constraint violation should be prevented by appropriate type checking in the service layer.
			// https://www.postgresql.org/docs/current/errcodes-appendix.html
			if strings.Contains(err.Error(), "SQLSTATE=23505") {
				return ErrUniqConstraintViolation
			}

			return ErrConstraintViolation
		} else if pgErr.StatementTimeout() {
			return fmt.Errorf("%w: %v", ErrTimeout, err)
		}
	}

	if strings.Contains(err.Error(), "AGORA=MISSINGRELATION") {
		return ErrMissingRelation
	}

	return err
}

func ForceRowsUpdate(res sql.Result) error {
	rows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check rows affected by the operation: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func NewErrInvalidEntity(field string, reason string) error {
	return fmt.Errorf("on field %q: %w: %s", field, ErrInvalidEntity, reason)
}

func NewErrNil(field string) error {
	return fmt.Errorf("on field %q: %w", field, ErrNil)
}

func NewErrInvalidCredentials(reason string) error {
	return fmt.Errorf("%w: %s", ErrInvalidCredentials, reason)
}

func NewErrUnauthorized(reason string) error {
	return fmt.Errorf("%w: %s", ErrUnauthorized, reason)
}

func NewErrNotAllowed[T any](field string, allowed ...T) error {
	return fmt.Errorf("on field %q: %w: allowed values are %v", field, ErrNotAllowed, allowed)
}

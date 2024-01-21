package db

import (
	"donation-mgmt/src/apperrors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type EntityIdentifier struct {
	EntityName string
	EntityID   string
	Extra      map[string]any
}

func MapDBError(err error, identifier EntityIdentifier) error {
	if err == nil {
		return nil
	}

	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		switch pgerr.Code {
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		case "23505": // unique violation
			return &apperrors.EntityAlreadyExistsError{
				EntityName: identifier.EntityName,
				EntityID:   identifier.EntityID,
				Extra:      identifier.Extra,
			}
		case "23503": // foreign key violation
			return &apperrors.EntityNotFoundError{
				EntityName: identifier.EntityName,
				EntityID:   identifier.EntityID,
				Extra:      identifier.Extra,
			}
		}
	}

	return fmt.Errorf("error executing database query: %w", err)
}

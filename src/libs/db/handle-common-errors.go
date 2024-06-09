package db

import (
	"donation-mgmt/src/apperrors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapDBError(err error, identifier apperrors.EntityIdentifier) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return &apperrors.EntityNotFoundError{
			EntityID: identifier,
		}
	}

	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		switch pgerr.Code {
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		case "23505": // unique violation
			return &apperrors.EntityAlreadyExistsError{
				EntityID: identifier,
			}
		case "23503": // foreign key violation
			return &apperrors.EntityNotFoundError{
				EntityID: identifier,
			}
		}
	}

	return fmt.Errorf("error executing %s DB Query: %w", identifier.EntityType, err)
}

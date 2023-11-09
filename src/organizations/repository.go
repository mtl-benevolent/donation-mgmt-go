package organizations

import (
	"context"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func FindOrganizationById(ctx context.Context, orgId uuid.UUID) (Organization, error) {
	dbQueries := data_access.New(db.DBPool())

	oid := &pgtype.UUID{}
	_ = oid.Scan(orgId.String())

	org, err := dbQueries.GetOrganizationByID(ctx, *oid)

	if err != nil {
		return Organization{}, err
	}

	model := Organization{
		ID:   orgId,
		Name: org.Name,
		Slug: org.Slug,

		// TODO: Deal with the pgtype correctly
		LogoUrl: nil,
	}

	return model, nil
}

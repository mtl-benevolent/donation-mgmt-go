package organizations

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"

	"github.com/jackc/pgx/v5"
)

type OrgService interface {
	GetOrganizationBySlug(ctx context.Context, slug string) (data_access.Organization, error)
	GetOrganizations(ctx context.Context) ([]data_access.Organization, error)
}

type OrgServiceImpl struct {
}

func NewOrgService() OrgService {
	return &OrgServiceImpl{}
}

func (s *OrgServiceImpl) GetOrganizationBySlug(ctx context.Context, slug string) (data_access.Organization, error) {
	// TODO: Implement Unit of Work pattern
	repo := data_access.New(db.DBPool())

	org, err := repo.GetOrganizationBySlug(ctx, slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return data_access.Organization{}, &apperrors.EntityNotFoundError{
				EntityName: "Organization",
				EntityID:   slug,
			}
		}

		return data_access.Organization{}, err
	}

	return org, nil
}

func (s *OrgServiceImpl) GetOrganizations(ctx context.Context) ([]data_access.Organization, error) {
	// TODO: Implement Unit of Work pattern
	repo := data_access.New(db.DBPool())

	orgs, err := repo.GetOrganizations(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []data_access.Organization{}, nil
		}

		return []data_access.Organization{}, err
	}

	return orgs, nil
}

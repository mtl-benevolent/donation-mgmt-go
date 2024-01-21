package organizations

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/pagination"

	"github.com/jackc/pgx/v5"
)

type OrgService interface {
	GetOrganizationBySlug(ctx context.Context, slug string) (data_access.Organization, error)
	GetOrganizations(ctx context.Context, page pagination.PaginationOptions) (pagination.PaginatedResult[data_access.Organization], error)
	CreateOrganization(ctx context.Context, params data_access.InsertOrganizationParams) (data_access.Organization, error)
	UpdateOrganization(ctx context.Context, params data_access.UpdateOrganizationBySlugParams) (data_access.Organization, error)
}

type OrgServiceImpl struct {
}

func NewOrgService() OrgService {
	return &OrgServiceImpl{}
}

func (s *OrgServiceImpl) GetOrganizationBySlug(ctx context.Context, slug string) (data_access.Organization, error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return data_access.Organization{}, err
	}

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

func (s *OrgServiceImpl) GetOrganizations(ctx context.Context, page pagination.PaginationOptions) (pagination.PaginatedResult[data_access.Organization], error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return pagination.PaginatedResult[data_access.Organization]{}, err
	}

	orgs, err := repo.GetOrganizations(ctx, data_access.GetOrganizationsParams{
		Offset: int32(page.Offset),
		Limit:  int32(page.Limit),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return pagination.PaginatedResult[data_access.Organization]{}, nil
		}

		return pagination.PaginatedResult[data_access.Organization]{}, err
	}

	total, err := repo.GetOrganizationsCount(ctx)
	if err != nil {
		return pagination.PaginatedResult[data_access.Organization]{}, err
	}

	paginatedResult := pagination.PaginatedResult[data_access.Organization]{
		Results: orgs,
		Total:   int(total),
	}

	return paginatedResult, nil
}

func (s *OrgServiceImpl) CreateOrganization(ctx context.Context, params data_access.InsertOrganizationParams) (data_access.Organization, error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return data_access.Organization{}, err
	}

	inserted, err := repo.InsertOrganization(ctx, params)
	if err != nil {
		return inserted, db.MapDBError(err, db.EntityIdentifier{
			EntityName: "Organization",
			EntityID:   params.Slug,
		})
	}

	return inserted, nil
}

func (s *OrgServiceImpl) UpdateOrganization(ctx context.Context, params data_access.UpdateOrganizationBySlugParams) (data_access.Organization, error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return data_access.Organization{}, err
	}

	updated, err := repo.UpdateOrganizationBySlug(ctx, params)
	if err != nil {
		return updated, db.MapDBError(err, db.EntityIdentifier{
			EntityName: "Organization",
			EntityID:   params.Slug,
		})
	}

	return updated, nil
}

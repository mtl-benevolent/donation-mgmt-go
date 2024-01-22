package organizations

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/pagination"

	"github.com/jackc/pgx/v5"
)

type OrganizationService struct {
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{}
}

func (s *OrganizationService) GetOrganizationBySlug(ctx context.Context, slug string) (data_access.Organization, error) {
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

type ListOrganizationsParams struct {
	// Subject is the user that is requesting the list of organizations. Leave empty to get all organizations.
	Subject     string
	PageOptions pagination.PaginationOptions
}

func (s *OrganizationService) GetOrganizations(ctx context.Context, params ListOrganizationsParams) (pagination.PaginatedResult[data_access.Organization], error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return pagination.PaginatedResult[data_access.Organization]{}, err
	}

	var orgs []data_access.Organization

	if params.Subject == "" {
		orgs, err = repo.ListOrganizations(ctx, data_access.ListOrganizationsParams{
			Offset: int32(params.PageOptions.Offset),
			Limit:  int32(params.PageOptions.Limit),
		})
	} else {
		orgs, err = repo.ListAuthorizedOrganizations(ctx, data_access.ListAuthorizedOrganizationsParams{
			Subject: params.Subject,
			Offset:  int32(params.PageOptions.Offset),
			Limit:   int32(params.PageOptions.Limit),
		})
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return pagination.PaginatedResult[data_access.Organization]{}, nil
		}

		return pagination.PaginatedResult[data_access.Organization]{}, err
	}

	var total int64

	if params.Subject == "" {
		total, err = repo.CountOrganizations(ctx)
	} else {
		total, err = repo.CountAuthorizedOrganizations(ctx, params.Subject)
	}

	if err != nil {
		return pagination.PaginatedResult[data_access.Organization]{}, err
	}

	paginatedResult := pagination.PaginatedResult[data_access.Organization]{
		Results: orgs,
		Total:   int(total),
	}

	return paginatedResult, nil
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, params data_access.InsertOrganizationParams) (data_access.Organization, error) {
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

func (s *OrganizationService) UpdateOrganization(ctx context.Context, params data_access.UpdateOrganizationBySlugParams) (data_access.Organization, error) {
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

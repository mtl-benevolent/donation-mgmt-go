package organizations

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/pagination"
	"errors"
	"fmt"

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
		return data_access.Organization{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			IDField:    "slug",
			EntityID:   slug,
		})
	}

	return org, nil
}

func (s *OrganizationService) GetOrganizationIDForSlug(ctx context.Context, slug string) (int64, error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return 0, err
	}

	orgID, err := repo.GetOrganizationIDBySlug(ctx, slug)
	if err != nil {
		return 0, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			IDField:    "slug",
			EntityID:   slug,
		})
	}

	return orgID, nil
}

func (s *OrganizationService) ListFiscalYearsForOrganization(ctx context.Context, orgID int64, environment data_access.Enviroment) ([]int16, error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return []int16{}, err
	}

	fiscalYears, err := repo.ListOrganizationFiscalYears(ctx, data_access.ListOrganizationFiscalYearsParams{
		OrganizationID: orgID,
		Environment:    environment,
	})

	if err != nil {
		return []int16{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			IDField:    "id",
			EntityID:   fmt.Sprintf("%d", orgID),
			Extras: map[string]interface{}{
				"environment": environment,
			},
		})
	}

	return fiscalYears, nil
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
		if errors.Is(err, pgx.ErrNoRows) {
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
		return inserted, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			Extras: map[string]interface{}{
				"name": params.Name,
				"slug": params.Slug,
			},
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
		return updated, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			EntityID:   params.Slug,
		})
	}

	return updated, nil
}

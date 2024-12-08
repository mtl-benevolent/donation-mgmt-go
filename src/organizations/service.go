package organizations

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
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

func (s *OrganizationService) GetOrganizationBySlug(ctx context.Context, querier dal.Querier, slug string) (dal.Organization, error) {
	org, err := querier.GetOrganizationBySlug(ctx, slug)
	if err != nil {
		return dal.Organization{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			IDField:    "slug",
			EntityID:   slug,
		})
	}

	return org, nil
}

func (s *OrganizationService) GetOrganizationIDForSlug(ctx context.Context, querier dal.Querier, slug string) (int64, error) {
	orgID, err := querier.GetOrganizationIDBySlug(ctx, slug)
	if err != nil {
		return 0, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			IDField:    "slug",
			EntityID:   slug,
		})
	}

	return orgID, nil
}

func (s *OrganizationService) ListFiscalYearsForOrganization(ctx context.Context, querier dal.Querier, orgID int64, environment dal.Enviroment) ([]int16, error) {
	fiscalYears, err := querier.ListOrganizationFiscalYears(ctx, dal.ListOrganizationFiscalYearsParams{
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

func (s *OrganizationService) GetOrganizations(ctx context.Context, querier dal.Querier, params ListOrganizationsParams) (pagination.PaginatedResult[dal.Organization], error) {
	var orgs []dal.Organization
	var err error

	if params.Subject == "" {
		orgs, err = querier.ListOrganizations(ctx, dal.ListOrganizationsParams{
			Offset: int32(params.PageOptions.Offset),
			Limit:  int32(params.PageOptions.Limit),
		})
	} else {
		orgs, err = querier.ListAuthorizedOrganizations(ctx, dal.ListAuthorizedOrganizationsParams{
			Subject: params.Subject,
			Offset:  int32(params.PageOptions.Offset),
			Limit:   int32(params.PageOptions.Limit),
		})
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pagination.PaginatedResult[dal.Organization]{}, nil
		}

		return pagination.PaginatedResult[dal.Organization]{}, err
	}

	var total int64

	if params.Subject == "" {
		total, err = querier.CountOrganizations(ctx)
	} else {
		total, err = querier.CountAuthorizedOrganizations(ctx, params.Subject)
	}

	if err != nil {
		return pagination.PaginatedResult[dal.Organization]{}, err
	}

	paginatedResult := pagination.PaginatedResult[dal.Organization]{
		Results: orgs,
		Total:   int(total),
	}

	return paginatedResult, nil
}

func (s *OrganizationService) CreateOrganization(ctx context.Context, querier dal.Querier, params dal.InsertOrganizationParams) (dal.Organization, error) {
	inserted, err := querier.InsertOrganization(ctx, params)
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

func (s *OrganizationService) UpdateOrganization(ctx context.Context, querier dal.Querier, params dal.UpdateOrganizationBySlugParams) (dal.Organization, error) {
	updated, err := querier.UpdateOrganizationBySlug(ctx, params)
	if err != nil {
		return updated, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			EntityID:   params.Slug,
		})
	}

	return updated, nil
}

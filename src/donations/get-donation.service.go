package donations

import (
	"context"
	"encoding/json"
	"fmt"

	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/libs/db"
)

type GetDonationByIDParams struct {
	OrganizationID int64
	Environment    dal.Environment
	DonationID     int64
}

func (s *DonationsService) GetDonationByID(ctx context.Context, querier dal.Querier, params GetDonationByIDParams) (DonationModel, error) {
	donationRows, err := querier.GetDonationByID(ctx, dal.GetDonationByIDParams{
		ID:             params.DonationID,
		OrganizationID: params.OrganizationID,
		Environment:    params.Environment,
	})

	if err != nil {
		return DonationModel{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Donation",
			IDField:    "id",
			EntityID:   fmt.Sprintf("%d", params.DonationID),
			Extras: map[string]interface{}{
				"organizationId": params.OrganizationID,
				"environment":    params.Environment,
			},
		})
	}

	return mapDonationRows(donationRows)
}

type GetDonationBySlugParams struct {
	OrganizationID int64
	Environment    dal.Environment
	Slug           string
}

func (s *DonationsService) GetDonationBySlug(ctx context.Context, querier dal.Querier, params GetDonationBySlugParams) (DonationModel, error) {
	donationRows, err := querier.GetDonationBySlug(ctx, dal.GetDonationBySlugParams{
		Slug:           params.Slug,
		OrganizationID: params.OrganizationID,
		Environment:    params.Environment,
	})

	if err != nil {
		return DonationModel{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Donation",
			IDField:    "slug",
			EntityID:   params.Slug,
			Extras: map[string]interface{}{
				"organizationId": params.OrganizationID,
				"environment":    params.Environment,
			},
		})
	}

	// Same struct, but we have to convert it to the correct struct array
	rows := make([]dal.GetDonationByIDRow, len(donationRows))
	for i, row := range donationRows {
		rows[i] = (dal.GetDonationByIDRow)(row)
	}

	return mapDonationRows(rows)
}

func mapDonationRows(donationRows []dal.GetDonationByIDRow) (DonationModel, error) {
	model := DonationModel{
		Payments: make([]dal.DonationPayment, len(donationRows)),
	}

	for i, row := range donationRows {
		if i == 0 {
			var donorAddr DonorAddress
			if err := json.Unmarshal(row.DonorAddress, &donorAddr); err != nil {
				return DonationModel{}, fmt.Errorf("failed to unmarshal donor address: %w", err)
			}

			model.ID = row.ID
			model.Slug = row.Slug
			model.OrganizationID = row.OrganizationID
			model.ExternalID = row.ExternalID
			model.Environment = row.Environment
			model.FiscalYear = row.FiscalYear
			model.Reason = row.Reason
			model.Type = row.Type
			model.Source = row.Source
			model.DonorFirstname = row.DonorFirstname
			model.DonorLastnameOrOrgName = row.DonorLastnameOrOrgName
			model.DonorEmail = row.DonorEmail
			model.DonorAddress = donorAddr
			model.Donation.DonorAddress = row.DonorAddress
			model.EmitReceipt = row.EmitReceipt
			model.SendByEmail = row.SendByEmail
			model.CreatedAt = row.CreatedAt
			model.UpdatedAt = row.UpdatedAt
			model.ArchivedAt = row.ArchivedAt
			model.CommentsCount = row.CommentsCount
		}

		model.Payments[i] = dal.DonationPayment{
			ID:                   row.ID_2,
			ExternalID:           row.ExternalID_2,
			DonationID:           row.DonationID,
			AmountInCents:        row.AmountInCents,
			ReceiptAmountInCents: row.ReceiptAmountInCents,
			ReceivedAt:           row.ReceivedAt,
			CreatedAt:            row.CreatedAt_2,
			ArchivedAt:           row.ArchivedAt_2,
		}
	}

	return model, nil
}

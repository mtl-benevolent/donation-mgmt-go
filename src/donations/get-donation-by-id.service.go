package donations

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *DonationsService) GetDonationByID(ctx context.Context, donationID int64, organizationID int64) (DonationModel, error) {
	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	querier, err := uow.GetQuerier(ctx)
	if err != nil {
		return DonationModel{}, err
	}

	donationRows, err := querier.GetDonationByID(ctx, data_access.GetDonationByIDParams{
		ID:             donationID,
		OrganizationID: organizationID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DonationModel{}, &apperrors.EntityNotFoundError{
				EntityName: "Donation",
				EntityID:   fmt.Sprintf("%d", donationID),
				Extra: map[string]interface{}{
					"organizationId": organizationID,
				},
			}
		}

		return DonationModel{}, db.MapDBError(err, db.EntityIdentifier{
			EntityName: "Donation",
			EntityID:   fmt.Sprintf("%d", donationID),
			Extra: map[string]interface{}{
				"organizationId": organizationID,
			},
		})
	}

	return mapDonationRowsToModel(donationRows)
}

func mapDonationRowsToModel(donationRows []data_access.GetDonationByIDRow) (DonationModel, error) {
	model := DonationModel{
		Payments: make([]data_access.DonationPayment, len(donationRows)),
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

		model.Payments[i] = data_access.DonationPayment{
			ID:            row.ID_2,
			ExternalID:    row.ExternalID_2,
			DonationID:    row.DonationID,
			Amount:        row.Amount,
			ReceiptAmount: row.ReceiptAmount,
			ReceivedAt:    row.ReceivedAt,
			CreatedAt:     row.CreatedAt_2,
			ArchivedAt:    row.ArchivedAt_2,
		}
	}

	return model, nil
}

package donations

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/system/logging"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

var errRecurrentDonationNotFound = errors.New("recurrent donation not found")

type CreateDonationParams struct {
	OrganizationID int64
	Environment    dal.Enviroment
	ExternalID     *string
	Reason         *string
	Type           dal.DonationType
	Source         dal.DonationSource

	DonorFirstName         *string
	DonorLastnameOrOrgName *string
	DonorEmail             *string
	DonorAddress           DonorAddress

	FiscalYear  *int16
	EmitReceipt bool
	SendByEmail bool

	PaymentAmountInCents int64
	ReceiptAmountInCents int64
	ReceivedAt           time.Time
	PaymentExternalID    *string
}

func (p CreateDonationParams) IsRecurrent() bool {
	return p.ExternalID != nil && p.Type == dal.DonationTypeRECURRENT
}

// AddPayment adds a payment to either an existing recurring donation or to a new donation. If no donation exists, a new one will be created.
// A payment can be added to a given donation if the donation is recurrent and if the ExternalID match
// an entry in the database. Otherwise, a new donation is created.
func (s *DonationsService) AddPayment(ctx context.Context, querier dal.Querier, params CreateDonationParams) (DonationModel, error) {
	l := logging.WithContextData(ctx, s.l)

	if params.FiscalYear == nil {
		l.Debug("Fiscal year not provided, extracting from received at", "received_at", params.ReceivedAt)

		org, err := s.orgSvc.GetOrganizationByID(ctx, querier, params.OrganizationID)
		if err != nil {
			return DonationModel{}, err
		}

		fiscalYear, err := extractFiscalYear(params.ReceivedAt, org.Timezone)
		if err != nil {
			return DonationModel{}, fmt.Errorf("failed to extract fiscal year: %w", err)
		}

		params.FiscalYear = &fiscalYear
	}

	if params.IsRecurrent() {
		l = l.With("external_id", params.ExternalID, "source", params.Source)

		l.Info("Donation payment is recurrent, trying to insert payment to existing donation")
		donation, err := s.tryInsertPayment(ctx, querier, dal.InsertPaymentToRecurrentDonationParams{
			ExternalID:           params.ExternalID,
			AmountInCents:        params.PaymentAmountInCents,
			ReceiptAmountInCents: params.ReceiptAmountInCents,
			FiscalYear:           *params.FiscalYear,
			OrganizationID:       params.OrganizationID,
			Environment:          params.Environment,
		})

		if err != nil {
			if !errors.Is(err, errRecurrentDonationNotFound) {
				return DonationModel{}, fmt.Errorf("failed to insert payment in donation: %w", err)
			}

			// We will create a new donation
			l.Info("Recurrent donation not found, creating new donation")
		} else {
			l.Info("Payment inserted to existing donation")
			return donation, nil
		}
	}

	l.Info("Creating new donation")
	insertDonation, err := mapParamsToInsertDonation(params)
	if err != nil {
		return DonationModel{}, fmt.Errorf("failed mapping donation to db model: %w", err)
	}

	insertPayment := mapParamsToInsertPayment(params)

	donation, err := s.insertDonation(ctx, querier, insertDonation, insertPayment)
	if err != nil {
		return DonationModel{}, fmt.Errorf("failed to insert donation: %w", err)
	}

	return donation, nil
}

func (s *DonationsService) tryInsertPayment(
	ctx context.Context,
	querier dal.Querier,
	payment dal.InsertPaymentToRecurrentDonationParams,
) (DonationModel, error) {
	inserted, err := querier.InsertPaymentToRecurrentDonation(ctx, payment)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DonationModel{}, errRecurrentDonationNotFound
		}

		return DonationModel{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "DonationPayment",
			Extras: map[string]interface{}{
				"externalID":     payment.ExternalID,
				"fiscalYear":     payment.FiscalYear,
				"organizationId": payment.OrganizationID,
			},
		})
	}

	return s.GetDonationByID(ctx, querier, GetDonationByIDParams{
		OrganizationID: payment.OrganizationID,
		Environment:    payment.Environment,
		DonationID:     inserted.DonationID,
	})
}

func mapParamsToInsertDonation(params CreateDonationParams) (dal.InsertDonationParams, error) {
	donorAddr, err := json.Marshal(params.DonorAddress)
	if err != nil {
		return dal.InsertDonationParams{}, fmt.Errorf("failed to marshal donor address: %w", err)
	}

	slug := ulid.Make().String()
	if params.FiscalYear == nil {
		return dal.InsertDonationParams{}, errors.New("fiscal year is required")
	}

	donationToInsert := dal.InsertDonationParams{
		Slug:                   slug,
		OrganizationID:         params.OrganizationID,
		ExternalID:             params.ExternalID,
		Environment:            params.Environment,
		FiscalYear:             *params.FiscalYear,
		Reason:                 params.Reason,
		Type:                   params.Type,
		Source:                 params.Source,
		DonorFirstname:         params.DonorFirstName,
		DonorLastNameOrOrgName: *params.DonorLastnameOrOrgName,
		DonorEmail:             params.DonorEmail,
		DonorAddress:           donorAddr,
		EmitReceipt:            params.EmitReceipt,
		SendByEmail:            params.SendByEmail,
	}

	return donationToInsert, nil
}

func extractFiscalYear(t time.Time, timezone string) (int16, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, fmt.Errorf("failed to load location %s: %w", timezone, err)
	}

	fiscalYear := t.In(location).Year()
	return int16(fiscalYear), nil
}

func mapParamsToInsertPayment(params CreateDonationParams) dal.InsertDonationPaymentParams {
	return dal.InsertDonationPaymentParams{
		DonationID:    -1,
		ExternalID:    params.PaymentExternalID,
		Amount:        params.PaymentAmountInCents,
		ReceiptAmount: params.ReceiptAmountInCents,
		ReceivedAt:    params.ReceivedAt,
	}
}

func (s *DonationsService) insertDonation(
	ctx context.Context,
	querier dal.Querier,
	donation dal.InsertDonationParams,
	payment dal.InsertDonationPaymentParams,
) (DonationModel, error) {
	insertedDonation, err := querier.InsertDonation(ctx, donation)
	if err != nil {
		return DonationModel{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Donation",
			Extras:     map[string]interface{}{"slug": donation.Slug, "externalID": donation.ExternalID},
		})
	}

	insertedPayment, err := querier.InsertDonationPayment(ctx, payment)
	if err != nil {
		return DonationModel{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "DonationPayment",
			Extras:     map[string]interface{}{"donationID": insertedDonation.ID, "externalID": payment.ExternalID},
		})
	}

	var donorAddr DonorAddress
	err = json.Unmarshal(insertedDonation.DonorAddress, &donorAddr)
	if err != nil {
		return DonationModel{}, fmt.Errorf("failed to unmarshal donor address: %w", err)
	}

	return DonationModel{
		Donation:     insertedDonation,
		DonorAddress: donorAddr,
		Payments: []dal.DonationPayment{
			insertedPayment,
		},
		CommentsCount: 0,
	}, nil
}

package donations

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/gin/ginext"
	"donation-mgmt/src/libs/gin/ginutils"
	"donation-mgmt/src/libs/gin/middlewares"
	"donation-mgmt/src/organizations"
	"donation-mgmt/src/permissions"
	"donation-mgmt/src/ptr"
	"donation-mgmt/src/system/contextual"
)

type ControllerV1 struct {
	donationsService *DonationsService
}

func NewControllerV1() *ControllerV1 {
	return &ControllerV1{
		donationsService: GetDonationsService(),
	}
}

func (c *ControllerV1) RegisterRoutes(router gin.IRouter) {
	group := router.Group(fmt.Sprintf("/v1/organizations/:%s/environments/:%s/donations", ginext.OrgSlugParamName, ginext.EnvParamName))

	readDonationPerm := permissions.Donation.Capability(permissions.Read)
	group.GET(fmt.Sprintf(":%s", ginext.DonationSlugParamName), middlewares.WithOrgAuthorization(ginext.OrgSlugParamName, readDonationPerm), c.GetDonationBySlugV1)

	createDonationPerm := permissions.Donation.Capability(permissions.Create)
	group.POST("", middlewares.WithOrgAuthorization(ginext.OrgSlugParamName, createDonationPerm), c.CreateDonationV1)
}

func (c *ControllerV1) GetDonationBySlugV1(ctx *gin.Context) {
	uow := db.NewUnitOfWork()
	defer uow.Finalize(ctx)

	querier, err := uow.GetQuerier(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	orgSlug := contextual.GetOrgSlug(ctx)
	orgID, err := organizations.GetOrgService().GetOrganizationIDForSlug(ctx, querier, orgSlug)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	env, err := contextual.GetValidEnv(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	slug := ctx.Params.ByName(ginext.DonationSlugParamName)
	if slug == "" {
		_ = ctx.Error(apperrors.NewInvalidParamError(ginext.DonationSlugParamName))
		return
	}

	donation, err := GetDonationsService().GetDonationBySlug(ctx, querier, GetDonationBySlugParams{
		OrganizationID: orgID,
		Environment:    env,
		Slug:           slug,
	})
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	dto := mapDonationToDTO(donation, false)
	ctx.JSON(http.StatusOK, dto)
}

func (c *ControllerV1) CreateDonationV1(ctx *gin.Context) {
	request, err := ginutils.DeserializeJSON[CreateDonationRequestV1](ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = ctx.Error(err)
		return
	}

	uow := db.NewUnitOfWorkWithTx()
	defer uow.Finalize(ctx)

	querier, err := uow.GetQuerier(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	orgSlug := contextual.GetOrgSlug(ctx)
	orgID, err := organizations.GetOrgService().GetOrganizationIDForSlug(ctx, querier, orgSlug)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	env, err := contextual.GetValidEnv(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	lastNameOrOrg := request.Donor.LastName
	if lastNameOrOrg == nil {
		lastNameOrOrg = request.Donor.OrgName
	}

	donation, err := GetDonationsService().AddPayment(ctx, querier, CreateDonationParams{
		OrganizationID: orgID,
		Environment:    env,
		Reason:         request.Reason,
		Source:         request.Source,

		DonorFirstName:         request.Donor.FirstName,
		DonorLastnameOrOrgName: lastNameOrOrg,
		DonorEmail:             request.Donor.Email,
		DonorAddress: DonorAddress{
			Line1:      request.Donor.Address.Line1,
			Line2:      request.Donor.Address.Line2,
			City:       request.Donor.Address.City,
			State:      request.Donor.Address.State,
			PostalCode: request.Donor.Address.PostalCode,
			Country:    request.Donor.Address.Country,
		},

		FiscalYear:  nil,
		EmitReceipt: request.EmitReceipt,
		SendByEmail: request.Donor.CommunicationChannel == CommunicationChannelEmail && ptr.UnwrapWithDefault(request.Donor.Email) != "",

		PaymentAmountInCents: request.AmountInCents,
		ReceiptAmountInCents: request.ReceiptAmountInCents,
		ReceivedAt:           request.ReceivedAt,

		// In this endpoint, we only allow the creation of one-time donations
		ExternalID:        nil,
		Type:              dal.DonationTypeONETIME,
		PaymentExternalID: nil,
	})

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if err = uow.Commit(ctx); err != nil {
		_ = ctx.Error(err)
		return
	}

	dto := mapDonationToDTO(donation, false)
	ctx.JSON(http.StatusCreated, dto)
}

func mapDonationToDTO(donation DonationModel, includeArchived bool) DonationDTO {
	dto := DonationDTO{
		ID:         donation.ID,
		Slug:       donation.Slug,
		Type:       donation.Type,
		FiscalYear: uint16(donation.FiscalYear),
		Reason:     ptr.Unwrap(donation.Reason, ""),
		Source:     donation.Source,
		ExternalID: donation.ExternalID,

		// Will calculate those when looping through payments
		TotalInCents:              0,
		TotalReceiptAmountInCents: 0,
		LastPaymentReceivedAt:     time.Time{},

		Payments: make([]PaymentDTO, 0, len(donation.Payments)),
		Donor: DonorDTO{
			Email: donation.DonorEmail,
			Address: &DonorAddressDTO{
				Line1:      donation.DonorAddress.Line1,
				Line2:      donation.DonorAddress.Line2,
				City:       donation.DonorAddress.City,
				State:      donation.DonorAddress.State,
				PostalCode: donation.DonorAddress.PostalCode,
				Country:    donation.DonorAddress.Country,
			},
			CommunicationChannel: CommunicationChannelSnailMail,
		},

		CreatedAt:  donation.CreatedAt,
		UpdatedAt:  donation.UpdatedAt,
		ArchivedAt: donation.ArchivedAt,
	}

	if donation.DonorFirstname == nil {
		dto.Donor.FirstName = donation.DonorFirstname
		dto.Donor.LastName = &donation.DonorLastnameOrOrgName
	} else {
		dto.Donor.OrgName = &donation.DonorLastnameOrOrgName
	}

	if donation.SendByEmail {
		dto.Donor.CommunicationChannel = CommunicationChannelEmail
	}

	for _, p := range donation.Payments {
		if p.ArchivedAt != nil && !includeArchived {
			continue
		}

		dto.Payments = append(dto.Payments, PaymentDTO{
			ID:                   p.ID,
			ExternalID:           p.ExternalID,
			AmountInCents:        p.AmountInCents,
			ReceiptAmountInCents: p.ReceiptAmountInCents,
			ReceivedAt:           p.ReceivedAt,
			CreatedAt:            p.CreatedAt,
			ArchivedAt:           p.ArchivedAt,
		})

		dto.TotalInCents += p.AmountInCents
		dto.TotalReceiptAmountInCents += p.ReceiptAmountInCents

		if dto.LastPaymentReceivedAt.Before(p.ReceivedAt) {
			dto.LastPaymentReceivedAt = p.ReceivedAt
		}
	}

	return dto
}

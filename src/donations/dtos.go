package donations

import (
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
	"reflect"
	"time"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CommunicationChannel string

const (
	CommunicationChannelEmail     CommunicationChannel = "EMAIL"
	CommunicationChannelSnailMail CommunicationChannel = "SNAIL_MAIL"
)

var validCommChannels = []any{
	CommunicationChannelEmail,
	CommunicationChannelSnailMail,
}

var validManualSources = []any{
	dal.DonationSourceCHEQUE,
	dal.DonationSourceDIRECTDEPOSIT,
	dal.DonationSourceOTHER,
	dal.DonationSourceSTOCKS,
}

type DonationDTO struct {
	ID                        int64              `json:"id"`
	Slug                      string             `json:"slug"`
	Type                      dal.DonationType   `json:"type"`
	FiscalYear                uint16             `json:"fiscalYear"`
	Reason                    string             `json:"reason,omitempty"`
	Source                    dal.DonationSource `json:"source"`
	ExternalID                *string            `json:"externalId,omitempty"`
	TotalInCents              int64              `json:"totalInCents"`
	TotalReceiptAmountInCents int64              `json:"totalReceiptAmountInCents"`
	LastPaymentReceivedAt     time.Time          `json:"lastPaymentReceivedAt"`

	Payments []PaymentDTO `json:"payments"`
	Donor    DonorDTO     `json:"donor"`

	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

type CreateDonationRequestV1 struct {
	Reason               *string            `json:"reason,omitempty"`
	Source               dal.DonationSource `json:"source"`
	AmountInCents        int64              `json:"amountInCents"`
	ReceiptAmountInCents int64              `json:"receiptAmountInCents"`
	ReceivedAt           time.Time          `json:"receivedAt"`

	Donor       DonorDTO `json:"donor"`
	EmitReceipt bool     `json:"emitReceipt"`
}

func (r CreateDonationRequestV1) Validate() error {
	err := ozzo.ValidateStruct(
		&r,
		ozzo.Field(&r.Reason, ozzo.Length(0, 255)),
		ozzo.Field(&r.Source, ozzo.Required, ozzo.In(validManualSources...)),
		ozzo.Field(&r.AmountInCents, ozzo.Required, ozzo.Min(1)),
		ozzo.Field(&r.ReceiptAmountInCents, ozzo.Required, ozzo.Min(1)),
		ozzo.Field(&r.ReceivedAt, ozzo.Required),
		ozzo.Field(&r.Donor, ozzo.NotNil),
	)

	if err != nil {
		return &apperrors.ValidationError{
			EntityName: reflect.TypeOf(r).Name(),
			InnerError: err,
		}
	}

	return nil
}

type DonorDTO struct {
	FirstName *string          `json:"firstName,omitempty"`
	LastName  *string          `json:"lastName,omitempty"`
	OrgName   *string          `json:"orgName,omitempty"`
	Email     *string          `json:"email,omitempty"`
	Address   *DonorAddressDTO `json:"address,omitempty"`

	CommunicationChannel CommunicationChannel `json:"communicationChannel"`
}

func (d DonorDTO) Validate() error {
	return ozzo.ValidateStruct(&d,
		ozzo.Field(&d.FirstName, ozzo.Length(0, 255)),
		ozzo.Field(&d.LastName, ozzo.Length(0, 255)),
		ozzo.Field(&d.OrgName, ozzo.Length(0, 255)),
		ozzo.Field(&d.Email, is.Email),
		ozzo.Field(&d.Address),
		ozzo.Field(&d.CommunicationChannel, ozzo.Required, ozzo.In(validCommChannels...)),
	)
}

type DonorAddressDTO struct {
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postalCode"`
	Country    *string `json:"country,omitempty"`
}

func (addr DonorAddressDTO) Validate() error {
	return ozzo.ValidateStruct(&addr,
		ozzo.Field(&addr.Line1, ozzo.Required, ozzo.Length(0, 255)),
		ozzo.Field(&addr.Line2, ozzo.Length(0, 255)),
		ozzo.Field(&addr.City, ozzo.Required, ozzo.Length(0, 255)),
		ozzo.Field(&addr.State, ozzo.Required, ozzo.Length(0, 255)),
		ozzo.Field(&addr.PostalCode, ozzo.Required, ozzo.Length(0, 255)),
		ozzo.Field(&addr.Country, ozzo.Length(0, 255)),
	)
}

type PaymentDTO struct {
	ID                   int64   `json:"id"`
	ExternalID           *string `json:"externalId"`
	AmountInCents        int64   `json:"amountInCents"`
	ReceiptAmountInCents int64   `json:"receiptAmountInCents"`

	ReceivedAt time.Time  `json:"receivedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

type CommentDTO struct {
	ID      int64  `json:"id"`
	Comment string `json:"comment"`
	Author  string `json:"author"`

	CreatedAt  time.Time  `json:"createdAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

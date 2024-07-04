package donations

import (
	"donation-mgmt/src/data_access"
	"time"

	ozzo "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var validSources = []any{
	data_access.DonationSourceCHEQUE,
	data_access.DonationSourceDIRECTDEPOSIT,
	data_access.DonationSourceOTHER,
	data_access.DonationSourceSTOCKS,
}

type CreateDonationRequestV1 struct {
	Reason               *string                    `json:"reason,omitempty"`
	Source               data_access.DonationSource `json:"source"`
	AmountInCents        int64                      `json:"amountInCents"`
	ReceiptAmountInCents int64                      `json:"receiptAmountInCents"`
	ReceivedAt           time.Time                  `json:"receivedAt"`

	Donor DonorDTO `json:"donor"`
}

func (r CreateDonationRequestV1) Validate() error {
	return ozzo.ValidateStruct(
		&r,
		ozzo.Field(&r.Reason, ozzo.Length(0, 255)),
		ozzo.Field(&r.Source, ozzo.In(validSources...)),
		ozzo.Field(&r.AmountInCents, ozzo.Required, ozzo.Min(1)),
		ozzo.Field(&r.ReceiptAmountInCents, ozzo.Required, ozzo.Min(1)),
		ozzo.Field(&r.ReceivedAt, ozzo.Required),
		ozzo.Field(&r.Donor, ozzo.NotNil),
	)
}

type DonorDTO struct {
	FirstName *string          `json:"firstName,omitempty"`
	LastName  *string          `json:"lastName,omitempty"`
	OrgName   *string          `json:"orgName,omitempty"`
	Email     *string          `json:"email,omitempty"`
	Address   *DonorAddressDTO `json:"address,omitempty"`

	EmitReceipt bool `json:"emitReceipt"`
	SendByEmail bool `json:"sendByEmail"`
}

func (d DonorDTO) Validate() error {
	return ozzo.ValidateStruct(&d,
		ozzo.Field(&d.FirstName, ozzo.Length(0, 255)),
		ozzo.Field(&d.LastName, ozzo.Length(0, 255)),
		ozzo.Field(&d.OrgName, ozzo.Length(0, 255)),
		ozzo.Field(&d.Email, ozzo.Required, is.Email),
		ozzo.Field(&d.Address),
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

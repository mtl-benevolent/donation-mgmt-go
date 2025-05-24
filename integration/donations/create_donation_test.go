package donations

import (
	"context"
	"donation-mgmt/integration"
	"donation-mgmt/integration/setup"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/donations"
	"donation-mgmt/src/ptr"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Smoke_CreateDonation_AsRoot(t *testing.T) {
	if ok := integration.Prepare(t); !ok {
		return
	}

	orgName := setup.GenerateName()
	orgSlug := setup.Slugify(orgName, 32)

	_ = setup.NewSetup().
		WithOrganization(orgName, orgSlug).
		Execute(context.Background(), t)

	donationDate := time.Now()
	fiscalYear := donationDate.Year()

	createDonationReq := donations.CreateDonationRequestV1{
		Reason:               nil,
		Source:               dal.DonationSourceCHEQUE,
		AmountInCents:        100_00,
		ReceiptAmountInCents: 50_00,
		ReceivedAt:           time.Now(),
		EmitReceipt:          true,
		Donor: donations.DonorDTO{
			FirstName:            ptr.Wrap("John"),
			LastName:             ptr.Wrap("Doe"),
			Email:                ptr.Wrap("john.doe@my-email.org"),
			CommunicationChannel: donations.CommunicationChannelSnailMail, // Prevents us from sending an email
			Address: &donations.DonorAddressDTO{
				Line1:      "22 Street Av.",
				City:       "Townsville",
				State:      "ON",
				PostalCode: "H0H 0H0",
				Country:    ptr.Wrap("CA"),
			},
		},
	}

	httpReq := setup.NewHttpReq(t, setup.HttpReqBuilder{
		Method: http.MethodPost,
		Url:    fmt.Sprintf("/v1/organizations/%s/environments/sandbox/donations", orgSlug),
		Body:   createDonationReq,
		User:   "root",
	})

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err, "Failed to make HTTP request")
	setup.AssertStatusCode(t, resp, http.StatusCreated)

	created, err := setup.ReadResponseBody[donations.DonationDTO](resp)
	require.NoError(t, err, "Failed to read response body")

	require.Equal(t, int64(100_00), created.TotalInCents, "Mismatching total amount")
	require.Equal(t, uint16(fiscalYear), created.FiscalYear, "Mismatching fiscal year")
	require.NotEmpty(t, created.Payments, "Donation should have payment")
	require.Equal(t, int64(50_00), created.Payments[0].ReceiptAmountInCents, "Mismatching receipt amount")
}

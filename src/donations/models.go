package donations

import "donation-mgmt/src/dal"

type DonationModel struct {
	dal.Donation

	DonorAddress DonorAddress

	CommentsCount int64
	Payments      []dal.DonationPayment
}

type DonorAddress struct {
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postalCode"`
	Country    *string `json:"country,omitempty"`
}

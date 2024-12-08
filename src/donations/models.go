package donations

import "donation-mgmt/src/dal"

type DonationModel struct {
	dal.Donation

	DonorAddress DonorAddress

	CommentsCount int64
	Payments      []dal.DonationPayment
}

type DonorAddress struct {
	Line1      string
	Line2      *string
	City       string
	State      string
	PostalCode string
	Country    *string
}

package donations

import "donation-mgmt/src/data_access"

type DonationModel struct {
	data_access.Donation

	DonorAddress DonorAddress

	CommentsCount int64
	Payments      []data_access.DonationPayment
}

type DonorAddress struct {
	Line1      string
	Line2      *string
	City       string
	State      string
	PostalCode string
	Country    *string
}

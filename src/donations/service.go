package donations

import (
	"log/slog"

	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/organizations"
)

type DonationsService struct {
	l      *slog.Logger
	orgSvc *organizations.OrganizationService
}

func NewDonationsService(orgSvc *organizations.OrganizationService) *DonationsService {
	return &DonationsService{
		l:      logger.ForComponent("donations-service"),
		orgSvc: orgSvc,
	}
}

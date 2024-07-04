package donations

import (
	"log/slog"

	"donation-mgmt/src/libs/logger"
)

type DonationsService struct {
	l *slog.Logger
}

func NewDonationsService() *DonationsService {
	return &DonationsService{
		l: logger.ForComponent("donations-service"),
	}
}

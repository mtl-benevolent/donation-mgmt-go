package donations

import (
	"donation-mgmt/src/organizations"

	"github.com/gin-gonic/gin"
)

var donationsService *DonationsService

func Bootstrap(router gin.IRouter) {
	donationsService = NewDonationsService(organizations.GetOrgService())

	v1 := NewControllerV1()
	v1.RegisterRoutes(router)
}

func GetDonationsService() *DonationsService {
	if donationsService == nil {
		panic("Donations service not bootstrapped")
	}

	return donationsService
}

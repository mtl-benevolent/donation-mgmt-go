package donations

import "github.com/gin-gonic/gin"

var donationsService *DonationsService

func Bootstrap(router gin.IRouter) {
	donationsService = NewDonationsService()

	v1 := NewControllerV1()
	v1.RegisterRoutes(router)
}

func GetDonationsService() *DonationsService {
	if donationsService == nil {
		panic("Donations service not bootstrapped")
	}

	return donationsService
}

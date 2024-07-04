package organizations

import "github.com/gin-gonic/gin"

var orgService *OrganizationService

func Bootstrap(router *gin.Engine) {
	if router != nil {
		registerRoutes(router)
	}

	orgService = NewOrganizationService()
}

func GetOrgService() *OrganizationService {
	if orgService == nil {
		panic("Organizations service not bootstrapped")
	}

	return orgService
}

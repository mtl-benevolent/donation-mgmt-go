package organizations

import "github.com/gin-gonic/gin"

var orgService *OrganizationService

func Bootstrap(router *gin.Engine) {
	orgService = NewOrganizationService()

	v1 := NewControllerV1()
	v1.RegisterRoutes(router)
}

func GetOrgService() *OrganizationService {
	if orgService == nil {
		panic("Organizations service not bootstrapped")
	}

	return orgService
}

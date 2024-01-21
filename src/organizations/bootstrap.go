package organizations

import "github.com/gin-gonic/gin"

var orgService OrgService

func Bootstrap(router *gin.Engine) {
	registerRoutes(router)

	orgService = NewOrgService()
}

func GetOrgService() OrgService {
	if orgService == nil {
		panic("Organizations service not bootstrapped")
	}

	return orgService
}

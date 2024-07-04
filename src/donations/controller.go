package donations

import (
	"donation-mgmt/src/libs/gin/ginutils"
	"donation-mgmt/src/libs/gin/middlewares"
	"donation-mgmt/src/permissions"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	orgSlugParam = "orgSlug"
)

type ControllerV1 struct {
	donationsService *DonationsService
}

func NewControllerV1() *ControllerV1 {
	return &ControllerV1{
		donationsService: GetDonationsService(),
	}
}

func (c *ControllerV1) RegisterRoutes(router gin.IRouter) {
	group := router.Group(fmt.Sprintf("/v1/organizations/:%s/donations", orgSlugParam))

	permissionCreate := permissions.Donation.Capability(permissions.Create)

	group.POST("", middlewares.WithOrgAuthorization(orgSlugParam, permissionCreate), c.CreateDonationV1)
}

func (c *ControllerV1) CreateDonationV1(ctx *gin.Context) {
	request, err := ginutils.DeserializeJSON[CreateDonationRequestV1](ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		ctx.Error(err)
		return
	}

}

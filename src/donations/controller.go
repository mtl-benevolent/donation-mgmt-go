package donations

import (
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/gin/ginutils"
	"donation-mgmt/src/libs/gin/middlewares"
	"donation-mgmt/src/organizations"
	"donation-mgmt/src/permissions"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	orgSlugParam = "orgSlug"
	envParam     = "env"
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
	group := router.Group(fmt.Sprintf("/v1/organizations/:%s/environments/%s/donations", orgSlugParam, envParam))

	permissionCreate := permissions.Donation.Capability(permissions.Create)

	group.POST("", middlewares.WithOrgAuthorization(orgSlugParam, permissionCreate), c.CreateDonationV1)
}

func (c *ControllerV1) CreateDonationV1(ctx *gin.Context) {
	request, err := ginutils.DeserializeJSON[CreateDonationRequestV1](ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := request.Validate(); err != nil {
		_ = ctx.Error(err)
		return
	}

	uow := db.NewUnitOfWorkWithTx()
	defer uow.Finalize(ctx)

	querier, err := uow.GetQuerier(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	orgSlug := ctx.Params.ByName(orgSlugParam)
	orgID, err := organizations.GetOrgService().GetOrganizationIDForSlug(ctx, querier, orgSlug)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	donation, err := GetDonationsService().AddPayment(ctx, querier, CreateDonationParams{
		OrganizationID: orgID,
	})
}

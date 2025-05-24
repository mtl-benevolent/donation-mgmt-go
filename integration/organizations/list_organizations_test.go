package organizations

import (
	"context"
	"donation-mgmt/integration"
	"donation-mgmt/integration/setup"
	"donation-mgmt/src/organizations"
	"donation-mgmt/src/pagination"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Smoke_ListOrganizations_AsRoot(t *testing.T) {
	if ok := integration.Prepare(t); !ok {
		return
	}

	orgName := setup.GenerateName()
	orgSlug := setup.Slugify(orgName, 32)

	_ = setup.NewSetup().
		WithOrganization(orgName, orgSlug).
		Execute(context.Background(), t)

	httpReq := setup.NewHttpReq(t, setup.HttpReqBuilder{
		Method: http.MethodGet,
		Url:    "/v1/organizations",
		User:   "root",
	})

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err, "Failed to make HTTP request")
	setup.AssertStatusCode(t, resp, http.StatusOK)

	list, err := setup.ReadResponseBody[pagination.PaginatedDTO[organizations.OrganizationDTOV1]](resp)
	require.NoError(t, err, "Failed to read response body")

	require.GreaterOrEqual(t, list.Total, 1, "Expected total to be at least 1")
	require.NotEmpty(t, list.Results, "Expected at least 1 result")
}

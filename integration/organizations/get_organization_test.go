package organizations

import (
	"context"
	"donation-mgmt/integration"
	"donation-mgmt/integration/setup"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/organizations"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Smoke_GetOrganizationBySlug_AsRoot(t *testing.T) {
	if ok := integration.Prepare(t); !ok {
		return
	}

	orgName := setup.GenerateName()
	orgSlug := setup.Slugify(orgName, 32)

	values := setup.NewSetup().
		WithOrganization(orgName, orgSlug).
		Execute(context.Background(), t)

	org, ok := setup.GetEntity[*dal.Organization](values, orgName)
	require.True(t, ok, "Failed to get organization entity")

	httpReq := setup.NewHttpReq(t, setup.HttpReqBuilder{
		Method: http.MethodGet,
		Url:    fmt.Sprintf("/v1/organizations/" + org.Slug),
		User:   "root",
	})

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err, "Failed to make HTTP request")
	setup.AssertStatusCode(t, resp, http.StatusOK)

	created, err := setup.ReadResponseBody[organizations.OrganizationDTOV1](resp)
	require.NoError(t, err, "Failed to read response body")

	require.Equal(t, org.Name, created.Name, "Name mismatch")
	require.Equal(t, org.Slug, created.Slug, "Slug mismatch")
}

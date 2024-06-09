package organizations

import (
	"donation-mgmt/integration"
	"donation-mgmt/integration/setup"
	"donation-mgmt/src/organizations"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Smoke_CreateOrganizationAsRoot(t *testing.T) {
	if ok := integration.Prepare(t); !ok {
		return
	}

	name := setup.GenerateName()

	req := organizations.CreateOrganizationRequest{
		Name: name,
		Slug: setup.Slugify(name, 32),
	}

	httpReq := setup.NewHttpReq(t, setup.HttpReqBuilder{
		Method: http.MethodPost,
		Url:    "/v1/organizations",
		Body:   req,
		User:   "root",
	})

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err, "Failed to make HTTP request")
	setup.AssertStatusCode(t, resp, http.StatusCreated)

	created, err := setup.ReadResponseBody[organizations.OrganizationDTO](resp)
	require.NoError(t, err, "Failed to read response body")

	require.Equal(t, req.Name, created.Name, "Name mismatch")
	require.Equal(t, req.Slug, created.Slug, "Slug mismatch")
}

package setup

import (
	"context"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/organizations"
	"errors"
	"fmt"
)

type OrganizationBuilder struct {
	name string
	slug string
}

func (s *Setup) WithOrganization(name string, slug string) *Setup {
	s.builders = append(s.builders, &OrganizationBuilder{
		name: name,
		slug: slug,
	})
	return s
}

func (b *OrganizationBuilder) Name() string {
	return b.name
}

func (b *OrganizationBuilder) Type() string {
	return "organization"
}

func (b *OrganizationBuilder) Execute(ctx context.Context, querier dal.Querier) (any, error) {
	if b.name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if b.slug == "" {
		b.slug = GenerateName()
	}

	org, err := organizations.GetOrgService().CreateOrganization(ctx, querier, dal.InsertOrganizationParams{
		Name: b.name,
		Slug: b.slug,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return &org, nil
}

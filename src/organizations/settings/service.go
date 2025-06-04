package settings

import (
	"context"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/encryption"
	"donation-mgmt/src/ptr"
	"encoding/json"
	"fmt"
)

type OrgSettingsService struct {
	encryptionKeyHex string
}

func NewOrgSettingsService(encryptionKeyHex string) *OrgSettingsService {
	return &OrgSettingsService{
		encryptionKeyHex: encryptionKeyHex,
	}
}

func (s *OrgSettingsService) InsertDefaultSettings(ctx context.Context, querier dal.Querier, orgID int64) error {
	envs := []dal.Environment{
		dal.EnvironmentSANDBOX,
		dal.EnvironmentLIVE,
	}

	for _, env := range envs {
		_, err := querier.UpsertOrganizationSettings(ctx, dal.UpsertOrganizationSettingsParams{
			OrganizationID: orgID,
			Environment:    env,
			EmailProvider: dal.NullEmailProvider{
				Valid: false,
			},
		})

		if err != nil {
			return fmt.Errorf("error inserting default settings for env %s: %w", env, err)
		}
	}

	return nil
}

type UpdateSettingsParams struct {
	OrgID                 int64
	Environment           dal.Environment
	Timezone              *string
	EmailProvider         *dal.EmailProvider
	EmailProviderSettings *EmailProviderSettings
}

func (s *OrgSettingsService) UpdateSettings(ctx context.Context, querier dal.Querier, params UpdateSettingsParams) (OrganizationSettings, error) {
	updates := dal.UpsertOrganizationSettingsParams{
		OrganizationID: params.OrgID,
		Environment:    params.Environment,
		Timezone:       params.Timezone,
	}

	if params.EmailProvider != nil {
		updates.EmailProvider = dal.NullEmailProvider{
			EmailProvider: *params.EmailProvider,
			Valid:         true,
		}
	}

	if params.EmailProviderSettings != nil {
		bytes, err := s.marshalSettings(*params.EmailProviderSettings)
		if err != nil {
			return OrganizationSettings{}, fmt.Errorf("error marshalling email provider settings: %w", err)
		}

		updates.EmailProviderSettings = bytes
	}

	updated, err := querier.UpsertOrganizationSettings(ctx, updates)
	if err != nil {
		return OrganizationSettings{}, fmt.Errorf("error updating organization settings: %w", err)
	}

	return s.MapDALToModel(updated)
}

func (s *OrgSettingsService) marshalSettings(settings EmailProviderSettings) ([]byte, error) {
	if settings.SMTP != nil && settings.SMTP.Password != nil {
		encrypted, err := encryption.EncryptString(*settings.SMTP.Password, s.encryptionKeyHex)
		if err != nil {
			return nil, fmt.Errorf("could not encrypt SMTP password")
		}

		settings.SMTP.Password = ptr.Wrap(encrypted)
	}

	return json.Marshal(settings)
}

func (s *OrgSettingsService) MapDALToModel(origin dal.OrganizationSetting) (OrganizationSettings, error) {

}

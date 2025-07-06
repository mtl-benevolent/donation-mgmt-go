package settings

import (
	"context"
	"donation-mgmt/src/apperrors"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/encryption"
	"donation-mgmt/src/libs/db"
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
		})

		if err != nil {
			return fmt.Errorf("error inserting default settings for env %s: %w", env, err)
		}
	}

	return nil
}

type GetSettingsParams struct {
	OrgID       int64
	Environment dal.Environment
}

type UpdateSettingsParams struct {
	OrgID       int64
	Environment dal.Environment
	Timezone    *string
}

type UpdateEmailSettingsParams struct {
	OrgID                 int64
	Environment           dal.Environment
	EmailProviderSettings EmailProviderSettings
}

func (s *OrgSettingsService) UpdateSettings(ctx context.Context, querier dal.Querier, params UpdateSettingsParams) (OrganizationSettings, error) {
	updates := dal.UpsertOrganizationSettingsParams{
		OrganizationID: params.OrgID,
		Environment:    params.Environment,
		Timezone:       params.Timezone,
	}

	updated, err := querier.UpsertOrganizationSettings(ctx, updates)
	if err != nil {
		return OrganizationSettings{}, db.MapDBError(err, apperrors.EntityIdentifier{
			EntityType: "Organization",
			IDField:    "OrgID",
			EntityID:   fmt.Sprintf("%d", params.OrgID),
			Extras: map[string]any{
				"Environment": params.Environment,
			},
		})
	}

	return s.MapDALToModel(updated)
}

func (s *OrgSettingsService) GetEmailSettings(ctx context.Context, querier dal.Querier, params GetSettingsParams) (EmailProviderSettings, error) {
	row, err := querier.GetOrganizationEmailSettings(ctx, dal.GetOrganizationEmailSettingsParams{
		OrganizationID: params.OrgID,
		Environment:    params.Environment,
	})

	errIdentitifer := apperrors.EntityIdentifier{
		EntityType: "Organization",
		IDField:    "OrgID",
		EntityID:   fmt.Sprintf("%d", params.OrgID),
		Extras: map[string]any{
			"Environment": params.Environment,
		},
	}

	if err != nil {
		return EmailProviderSettings{}, db.MapDBError(err, errIdentitifer)
	}

	if len(row.EmailProviderSettings) <= 0 {
		return EmailProviderSettings{}, &apperrors.EntityNotFoundError{
			EntityID: errIdentitifer,
		}
	}

	var encryptedSettings EncryptedEmailProviderSettings
	if err := json.Unmarshal([]byte(row.EmailProviderSettings), &encryptedSettings); err != nil {
		return EmailProviderSettings{}, fmt.Errorf("error unmarshaling email settings: %w", err)
	}

	settings := EmailProviderSettings{
		Provider: encryptedSettings.Provider,
	}

	// SMTP Provider
	if settings.Provider == SMTPEmailProvider && encryptedSettings.EncryptedSMTP != nil {
		smtp, err := encryption.DecryptJSON[SMTPSettings](*encryptedSettings.EncryptedSMTP, s.encryptionKeyHex)
		if err != nil {
			return EmailProviderSettings{}, fmt.Errorf("unable to decrypt email settings: %w", err)
		}

		settings.SMTP = &smtp
	}

	return settings, nil
}

func (s *OrgSettingsService) UpdateEmailSettings(ctx context.Context, querier dal.Querier, params UpdateEmailSettingsParams) error {
	encryptedVal, err := encryption.EncryptJSON(params.EmailProviderSettings, s.encryptionKeyHex)
	if err != nil {
		return fmt.Errorf("unable to encrypt email settings: %w", err)
	}

	updatedCount, err := querier.UpdateOrganizationEmailSettings(ctx, dal.UpdateOrganizationEmailSettingsParams{
		OrganizationID:        params.OrgID,
		Environment:           params.Environment,
		EmailProviderSettings: encryptedVal,
	})

	entityID := apperrors.EntityIdentifier{
		EntityType: "EmailSettings",
		IDField:    "OrganizationID",
		EntityID:   fmt.Sprintf("%d", params.OrgID),
		Extras: map[string]any{
			"Environment": params.Environment,
		},
	}

	if err != nil {
		return db.MapDBError(err, entityID)
	}

	if updatedCount == 0 {
		return &apperrors.EntityNotFoundError{
			EntityID: entityID,
		}
	}

	return nil
}

func (s *OrgSettingsService) MapDALToModel(origin dal.OrganizationSetting) (OrganizationSettings, error) {
	model := OrganizationSettings{
		OrganizationID: origin.OrganizationID,
		Environment:    origin.Environment,
		Timezone:       origin.Timezone,
		IsValid:        origin.IsValid,
		UpdatedAt:      origin.UpdatedAt,
	}

	return model, nil
}

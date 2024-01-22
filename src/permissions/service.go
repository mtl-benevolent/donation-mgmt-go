package permissions

import (
	"context"
	"donation-mgmt/src/data_access"
	"donation-mgmt/src/libs/db"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/system/contextual"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

var ErrInvalidParams = errors.New("invalid parameters")

type PermissionsService struct {
	l *slog.Logger
}

func NewPermissionsService() *PermissionsService {
	return &PermissionsService{
		l: logger.ForComponent("PermissionsService"),
	}
}

type HasRequiredPermissionsParams struct {
	OrganizationSlug string
	OrganizationID   int64

	MustBeGlobal bool

	Subject      string
	Capabilities []string
}

func (s *PermissionsService) HasCapabilities(ctx context.Context, params HasRequiredPermissionsParams) (bool, error) {
	if len(params.Capabilities) == 0 {
		return true, nil
	}

	uow, finalizer := db.GetUnitOfWorkFromCtxOrDefault(ctx)
	defer finalizer()

	repo, err := uow.GetQuerier(ctx)
	if err != nil {
		return false, err
	}

	if params.MustBeGlobal {
		return s.checkGlobal(ctx, repo, params)
	} else if params.OrganizationSlug != "" {
		return s.checkForOrgBySlug(ctx, repo, params)
	} else if params.OrganizationID > 0 {
		return s.checkForOrgByID(ctx, repo, params)
	} else {
		return false, ErrInvalidParams
	}
}

func (s *PermissionsService) checkGlobal(ctx context.Context, querier data_access.Querier, params HasRequiredPermissionsParams) (bool, error) {
	logger := contextual.LoggerWithContextData(ctx, s.l)

	role, err := querier.HasGlobalCapabilities(ctx, data_access.HasGlobalCapabilitiesParams{
		Capabilities: params.Capabilities,
		Subject:      params.Subject,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("User does not have required global permissions", slog.String("subject", params.Subject), slog.Any("permissions", params.Capabilities))
			return false, nil
		}

		logger.Error("Error checking global permissions", slog.Any("error", err))
		return false, err
	}

	logger.Debug("User has required global permissions", slog.String("subject", params.Subject), slog.Any("permissions", params.Capabilities), slog.String("role", role.Name))
	return true, nil
}

func (s *PermissionsService) checkForOrgBySlug(ctx context.Context, querier data_access.Querier, params HasRequiredPermissionsParams) (bool, error) {
	logger := contextual.LoggerWithContextData(ctx, s.l)

	role, err := querier.HasCapabilitiesForOrgBySlug(ctx, data_access.HasCapabilitiesForOrgBySlugParams{
		Capabilities:     params.Capabilities,
		Subject:          params.Subject,
		OrganizationSlug: params.OrganizationSlug,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("User does not have required permissions for org", slog.String("subject", params.Subject), slog.Any("permissions", params.Capabilities), slog.String("orgSlug", params.OrganizationSlug))
			return false, nil
		}

		logger.Error("Error checking permissions for org", slog.Any("error", err))
		return false, err
	}

	logger.Debug("User has required permissions for org", slog.String("subject", params.Subject), slog.Any("permissions", params.Capabilities), slog.String("orgSlug", params.OrganizationSlug), slog.String("role", role.Name))
	return true, nil
}

func (s *PermissionsService) checkForOrgByID(ctx context.Context, querier data_access.Querier, params HasRequiredPermissionsParams) (bool, error) {
	logger := contextual.LoggerWithContextData(ctx, s.l)

	role, err := querier.HasCapabilitiesForOrgByID(ctx, data_access.HasCapabilitiesForOrgByIDParams{
		Capabilities:   params.Capabilities,
		Subject:        params.Subject,
		OrganizationID: params.OrganizationID,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Debug("User does not have required permissions for org", slog.String("subject", params.Subject), slog.Any("permissions", params.Capabilities), slog.Int64("orgID", params.OrganizationID))
			return false, nil
		}

		logger.Error("Error checking permissions for org", slog.Any("error", err))
		return false, err
	}

	logger.Debug("User has required permissions for org", slog.String("subject", params.Subject), slog.Any("permissions", params.Capabilities), slog.Int64("orgID", params.OrganizationID), slog.String("role", role.Name))
	return true, nil
}

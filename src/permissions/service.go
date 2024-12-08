package permissions

import (
	"context"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/system/logging"
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

func (s *PermissionsService) HasCapabilities(ctx context.Context, querier dal.Querier, params HasRequiredPermissionsParams) (bool, error) {
	if len(params.Capabilities) == 0 {
		return true, nil
	}

	if params.MustBeGlobal {
		return s.checkGlobal(ctx, querier, params)
	} else if params.OrganizationSlug != "" {
		return s.checkForOrgBySlug(ctx, querier, params)
	} else if params.OrganizationID > 0 {
		return s.checkForOrgByID(ctx, querier, params)
	} else {
		return false, ErrInvalidParams
	}
}

func (s *PermissionsService) checkGlobal(ctx context.Context, querier dal.Querier, params HasRequiredPermissionsParams) (bool, error) {
	logger := logging.WithContextData(ctx, s.l)

	role, err := querier.HasGlobalCapabilities(ctx, dal.HasGlobalCapabilitiesParams{
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

func (s *PermissionsService) checkForOrgBySlug(ctx context.Context, querier dal.Querier, params HasRequiredPermissionsParams) (bool, error) {
	logger := logging.WithContextData(ctx, s.l)

	role, err := querier.HasCapabilitiesForOrgBySlug(ctx, dal.HasCapabilitiesForOrgBySlugParams{
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

func (s *PermissionsService) checkForOrgByID(ctx context.Context, querier dal.Querier, params HasRequiredPermissionsParams) (bool, error) {
	logger := logging.WithContextData(ctx, s.l)

	role, err := querier.HasCapabilitiesForOrgByID(ctx, dal.HasCapabilitiesForOrgByIDParams{
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

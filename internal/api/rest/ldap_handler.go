package rest

import (
	"context"
	"fmt"
	"time"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) CancelLDAPSync(ctx context.Context) (generatedapi.CancelLDAPSyncRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	if err := r.ldapUseCase.CancelSync(ctx); err != nil {
		return &generatedapi.Error{
			Error: generatedapi.ErrorError{
				Message: generatedapi.NewOptString(fmt.Sprintf("Failed to cancel sync: %v", err)),
			},
		}, nil
	}

	return &generatedapi.SuccessResponse{
		Message: generatedapi.NewOptString("Synchronization cancelled successfully"),
	}, nil
}

func (r *RestAPI) DeleteLDAPConfig(ctx context.Context) (generatedapi.DeleteLDAPConfigRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}

	return &generatedapi.SuccessResponse{
		Message: generatedapi.NewOptString("LDAP configuration deleted successfully"),
	}, nil
}

func (r *RestAPI) GetLDAPConfig(ctx context.Context) (generatedapi.GetLDAPConfigRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}

	config, err := r.settingsUseCase.GetLDAPConfig(ctx)
	if err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.LDAPConfig{
		Enabled:       config.Enabled,
		URL:           config.URL,
		BindDn:        config.BindDN,
		BindPassword:  config.BindPassword,
		UserBaseDn:    config.UserBaseDN,
		UserFilter:    config.UserFilter,
		UserNameAttr:  config.UserNameAttr,
		UserEmailAttr: config.UserEmailAttr,
		StartTLS:      config.StartTLS,
		InsecureTLS:   config.InsecureTLS,
		Timeout:       config.Timeout,
		SyncInterval:  config.SyncInterval,
	}, nil
}

func (r *RestAPI) GetLDAPStatistics(ctx context.Context) (generatedapi.GetLDAPStatisticsRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	stats, err := r.ldapUseCase.GetStatistics(ctx)
	if err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.LDAPStatistics{
		LdapUsers:       generatedapi.NewOptInt(stats.LDAPUsers),
		LocalUsers:      generatedapi.NewOptInt(stats.LocalUsers),
		ActiveUsers:     generatedapi.NewOptInt(stats.ActiveUsers),
		InactiveUsers:   generatedapi.NewOptInt(stats.InactiveUsers),
		SyncSuccessRate: generatedapi.NewOptFloat32(stats.SyncSuccessRate),
	}, nil
}

func (r *RestAPI) GetLDAPSyncLogDetails(
	ctx context.Context,
	params generatedapi.GetLDAPSyncLogDetailsParams,
) (generatedapi.GetLDAPSyncLogDetailsRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	log, err := r.ldapUseCase.GetSyncLogDetails(ctx, params.ID)
	if err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.LDAPSyncLogDetails{
		ID:               log.ID,
		Timestamp:        log.Timestamp,
		Level:            generatedapi.LDAPSyncLogDetailsLevel(log.Level),
		Message:          log.Message,
		SyncSessionID:    log.SyncSessionID,
		Username:         convertNilString(log.Username),
		Details:          convertNilString(log.Details),
		StackTrace:       convertNilString(log.StackTrace),
		LdapErrorCode:    convertNilInt(log.LDAPErrorCode),
		LdapErrorMessage: convertNilString(log.LDAPErrorMessage),
	}, nil
}

func (r *RestAPI) GetLDAPSyncLogs(
	ctx context.Context,
	params generatedapi.GetLDAPSyncLogsParams,
) (generatedapi.GetLDAPSyncLogsRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	filter := domain.LDAPSyncLogFilter{
		Limit:    convertOptInt(params.Limit),
		Level:    convertOptStringFromLevel(params.Level),
		SyncID:   convertOptString(params.SyncID),
		Username: convertOptString(params.Username),
		From:     convertOptDateTime(params.From),
		To:       convertOptDateTime(params.To),
	}
	result, err := r.ldapUseCase.GetSyncLogs(ctx, filter)
	if err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}
	logs := make([]generatedapi.LDAPSyncLogEntry, 0, len(result.Logs))
	for _, syncLog := range result.Logs {
		logs = append(logs, generatedapi.LDAPSyncLogEntry{
			ID:            syncLog.ID,
			Timestamp:     syncLog.Timestamp,
			Level:         generatedapi.LDAPSyncLogEntryLevel(syncLog.Level),
			Message:       syncLog.Message,
			SyncSessionID: syncLog.SyncSessionID,
			Username:      convertNilString(syncLog.Username),
			Details:       convertNilString(syncLog.Details),
		})
	}

	return &generatedapi.LDAPSyncLogs{
		Logs:    logs,
		Total:   generatedapi.NewOptInt(result.Total),
		HasMore: generatedapi.NewOptBool(result.HasMore),
	}, nil
}

func (r *RestAPI) GetLDAPSyncProgress(ctx context.Context) (generatedapi.GetLDAPSyncProgressRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	progress, err := r.ldapUseCase.GetSyncProgress(ctx)
	if err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.LDAPSyncProgress{
		IsRunning:      progress.IsRunning,
		Progress:       float32(progress.Progress),
		CurrentStep:    progress.CurrentStep,
		ProcessedItems: progress.ProcessedItems,
		TotalItems:     progress.TotalItems,
		EstimatedTime:  progress.EstimatedTime,
		StartTime:      progress.StartTime,
		SyncID:         progress.SyncID,
	}, nil
}

func (r *RestAPI) GetLDAPSyncStatus(ctx context.Context) (generatedapi.GetLDAPSyncStatusRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	status, err := r.ldapUseCase.GetSyncStatus(ctx)
	if err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.LDAPSyncStatus{
		Status:           status.Status,
		IsRunning:        status.IsRunning,
		LastSyncTime:     generatedapi.NewOptDateTime(status.LastSyncTime),
		TotalUsers:       status.TotalUsers,
		SyncedUsers:      status.SyncedUsers,
		Errors:           status.Errors,
		Warnings:         status.Warnings,
		LastSyncDuration: generatedapi.NewOptString(status.LastSyncDuration),
	}, nil
}

func (r *RestAPI) SyncLDAPUsers(ctx context.Context) (generatedapi.SyncLDAPUsersRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	if err := r.ldapUseCase.StartManualSync(ctx); err != nil {
		return &generatedapi.Error{
			Error: generatedapi.ErrorError{
				Message: generatedapi.NewOptString(fmt.Sprintf("Failed to start sync: %v", err)),
			},
		}, nil
	}

	return &generatedapi.LDAPSyncStartResponse{
		Message: generatedapi.NewOptString("LDAP users synchronization started successfully"),
	}, nil
}

func (r *RestAPI) TestLDAPConnection(
	ctx context.Context, _ *generatedapi.LDAPConnectionTest,
) (generatedapi.TestLDAPConnectionRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}
	if err := r.ldapUseCase.TestConnection(ctx); err != nil {
		return &generatedapi.LDAPConnectionTestResponse{
			Success: generatedapi.NewOptBool(false),
			Message: generatedapi.NewOptString(fmt.Sprintf("Connection test failed: %v", err)),
		}, nil
	}

	return &generatedapi.LDAPConnectionTestResponse{
		Success: generatedapi.NewOptBool(true),
		Message: generatedapi.NewOptString("LDAP connection test successful"),
	}, nil
}

func (r *RestAPI) UpdateLDAPConfig(
	ctx context.Context,
	req *generatedapi.LDAPConfig,
) (generatedapi.UpdateLDAPConfigRes, error) {
	userID := etogglcontext.UserID(ctx)
	if userID == 0 {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("Unauthorized"),
		}}, nil
	}
	if !etogglcontext.IsSuper(ctx) {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("Superuser access required"),
		}}, nil
	}

	config := &domain.LDAPConfig{
		Enabled:       req.Enabled,
		URL:           req.URL,
		BindDN:        req.BindDn,
		BindPassword:  req.BindPassword,
		UserBaseDN:    req.UserBaseDn,
		UserFilter:    req.UserFilter,
		UserNameAttr:  req.UserNameAttr,
		UserEmailAttr: req.UserEmailAttr,
		StartTLS:      req.StartTLS,
		InsecureTLS:   req.InsecureTLS,
		Timeout:       req.Timeout,
		SyncInterval:  req.SyncInterval,
	}

	err := r.ldapUseCase.UpdateConfig(ctx, config)
	if err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString(err.Error()),
		}}, nil
	}

	return &generatedapi.LDAPConfigResponse{
		Message: generatedapi.NewOptString("LDAP configuration updated successfully"),
		Config: generatedapi.NewOptLDAPConfig(generatedapi.LDAPConfig{
			Enabled:       config.Enabled,
			URL:           config.URL,
			BindDn:        config.BindDN,
			BindPassword:  config.BindPassword,
			UserBaseDn:    config.UserBaseDN,
			UserFilter:    config.UserFilter,
			UserNameAttr:  config.UserNameAttr,
			UserEmailAttr: config.UserEmailAttr,
			StartTLS:      config.StartTLS,
			InsecureTLS:   config.InsecureTLS,
			Timeout:       config.Timeout,
			SyncInterval:  config.SyncInterval,
		}),
	}, nil
}

// Helper functions for converting between domain and API types

func convertNilString(s *string) generatedapi.OptNilString {
	if s == nil {
		return generatedapi.OptNilString{}
	}

	return generatedapi.NewOptNilString(*s)
}

func convertNilInt(i *int) generatedapi.OptNilInt {
	if i == nil {
		return generatedapi.OptNilInt{}
	}

	return generatedapi.NewOptNilInt(*i)
}

func convertOptInt(opt generatedapi.OptInt) *int {
	if !opt.Set {
		return nil
	}

	return &opt.Value
}

func convertOptString(opt generatedapi.OptString) *string {
	if !opt.Set {
		return nil
	}

	return &opt.Value
}

func convertOptDateTime(opt generatedapi.OptDateTime) *time.Time {
	if !opt.Set {
		return nil
	}

	return &opt.Value
}

func convertOptStringFromLevel(opt generatedapi.OptGetLDAPSyncLogsLevel) *string {
	if !opt.Set {
		return nil
	}
	levelStr := string(opt.Value)

	return &levelStr
}

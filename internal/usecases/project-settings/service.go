package project_settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
	simplecache "github.com/togglr-project/togglr/pkg/simple-cache"
)

const (
	AutoDisableEnabledKey          = "auto_disable_enabled"
	AutoDisableRequiresApprovalKey = "auto_disable_requires_approval"
	AutoDisableErrorThresholdKey   = "auto_disable_error_threshold"
	AutoDisableTimeWindowSecKey    = "auto_disable_time_window_sec"

	DefaultAutoDisableEnabled = true
	DefaultRequiresApproval   = false
	DefaultErrorThreshold     = 10
	DefaultTimeWindow         = 60 * time.Second

	CacheTTL = 5 * time.Minute
)

type Service struct {
	txManager           db.TxManager
	projectSettingsRepo contract.ProjectSettingsRepository
	cache               *simplecache.Cache[string, any]
}

func New(txManager db.TxManager, projectSettingsRepo contract.ProjectSettingsRepository) *Service {
	service := &Service{
		txManager:           txManager,
		projectSettingsRepo: projectSettingsRepo,
		cache:               simplecache.New[string, any](),
	}
	service.cache.StartCleanup(1 * time.Minute)

	return service
}

func (s *Service) Create(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	value any,
) (*domain.ProjectSetting, error) {
	if name == "" {
		return nil, errors.New("setting name cannot be empty")
	}

	_, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("invalid setting value: %w", err)
	}

	setting := &domain.ProjectSetting{
		ProjectID: projectID,
		Name:      name,
		Value:     value,
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.projectSettingsRepo.Create(ctx, setting)
	})
	if err != nil {
		return nil, fmt.Errorf("create project setting: %w", err)
	}

	return setting, nil
}

func (s *Service) GetByName(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
) (*domain.ProjectSetting, error) {
	if name == "" {
		return nil, errors.New("setting name cannot be empty")
	}

	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("get project setting: %w", err)
	}

	return setting, nil
}

func (s *Service) Update(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	value any,
) (*domain.ProjectSetting, error) {
	if name == "" {
		return nil, errors.New("setting name cannot be empty")
	}

	_, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("invalid setting value: %w", err)
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.projectSettingsRepo.Update(ctx, projectID, name, value)
	})
	if err != nil {
		return nil, fmt.Errorf("update project setting: %w", err)
	}

	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil {
		return nil, fmt.Errorf("get updated project setting: %w", err)
	}

	return setting, nil
}

func (s *Service) Delete(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
) error {
	if name == "" {
		return errors.New("setting name cannot be empty")
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.projectSettingsRepo.Delete(ctx, projectID, name)
	})
	if err != nil {
		return fmt.Errorf("delete project setting: %w", err)
	}

	return nil
}

func (s *Service) List(
	ctx context.Context,
	projectID domain.ProjectID,
	page, perPage int,
) ([]*domain.ProjectSetting, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	settings, total, err := s.projectSettingsRepo.List(ctx, projectID, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("list project settings: %w", err)
	}

	return settings, total, nil
}

func (s *Service) GetBoolSetting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue bool,
) (bool, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case bool:
		return val, nil
	case string:
		if parsed, err := strconv.ParseBool(val); err == nil {
			return parsed, nil
		}

		return defaultValue, nil
	default:
		return defaultValue, nil
	}
}

func (s *Service) GetIntSetting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue int,
) (int, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed, nil
		}

		return defaultValue, nil
	default:
		return defaultValue, nil
	}
}

func (s *Service) GetStringSetting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue string,
) (string, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case string:
		return val, nil
	case bool:
		return strconv.FormatBool(val), nil
	case int:
		return strconv.Itoa(val), nil
	case int8:
		return strconv.Itoa(int(val)), nil
	case int16:
		return strconv.Itoa(int(val)), nil
	case int32:
		return strconv.Itoa(int(val)), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case uint:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(val), 10), nil
	case uint64:
		return strconv.FormatUint(val, 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	default:
		return defaultValue, nil
	}
}

func (s *Service) GetFloat64Setting(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	defaultValue float64,
) (float64, error) {
	setting, err := s.projectSettingsRepo.GetByName(ctx, projectID, name)
	if err != nil || setting == nil {
		return defaultValue, nil
	}

	switch val := setting.Value.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case string:
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed, nil
		}

		return defaultValue, nil
	default:
		return defaultValue, nil
	}
}

func (s *Service) GetAutoDisableEnabled(ctx context.Context, projectID domain.ProjectID) (bool, error) {
	return s.GetBoolSetting(ctx, projectID, AutoDisableEnabledKey, DefaultAutoDisableEnabled)
}

func (s *Service) GetAutoDisableRequiresApproval(ctx context.Context, projectID domain.ProjectID) (bool, error) {
	return s.GetBoolSetting(ctx, projectID, AutoDisableRequiresApprovalKey, DefaultRequiresApproval)
}

func (s *Service) GetAutoDisableErrorThreshold(ctx context.Context, projectID domain.ProjectID) (int, error) {
	return s.GetIntSetting(ctx, projectID, AutoDisableErrorThresholdKey, DefaultErrorThreshold)
}

func (s *Service) GetAutoDisableTimeWindowSec(ctx context.Context, projectID domain.ProjectID) (int, error) {
	return s.GetIntSetting(ctx, projectID, AutoDisableTimeWindowSecKey, int(DefaultTimeWindow.Seconds()))
}

func (s *Service) GetAutoDisableTimeWindow(ctx context.Context, projectID domain.ProjectID) (time.Duration, error) {
	secs, err := s.GetAutoDisableTimeWindowSec(ctx, projectID)
	if err != nil {
		return DefaultTimeWindow, err
	}

	return time.Duration(secs) * time.Second, nil
}

func (s *Service) GetAutoDisableEnabledCached(ctx context.Context, projectID domain.ProjectID) (bool, error) {
	cacheKey := makeProjectSettingCacheKey(projectID, AutoDisableEnabledKey)

	if cached, found := s.cache.Get(cacheKey); found {
		if value, ok := cached.(bool); ok {
			return value, nil
		}
	}

	value, err := s.GetAutoDisableEnabled(ctx, projectID)
	if err != nil {
		return false, err
	}

	s.cache.Set(cacheKey, value, CacheTTL)

	return value, nil
}

func (s *Service) GetAutoDisableRequiresApprovalCached(ctx context.Context, projectID domain.ProjectID) (bool, error) {
	cacheKey := makeProjectSettingCacheKey(projectID, AutoDisableRequiresApprovalKey)

	if cached, found := s.cache.Get(cacheKey); found {
		if value, ok := cached.(bool); ok {
			return value, nil
		}
	}

	value, err := s.GetAutoDisableRequiresApproval(ctx, projectID)
	if err != nil {
		return false, err
	}

	s.cache.Set(cacheKey, value, CacheTTL)

	return value, nil
}

func (s *Service) GetAutoDisableErrorThresholdCached(ctx context.Context, projectID domain.ProjectID) (int, error) {
	cacheKey := makeProjectSettingCacheKey(projectID, AutoDisableErrorThresholdKey)

	if cached, found := s.cache.Get(cacheKey); found {
		if value, ok := cached.(int); ok {
			return value, nil
		}
	}

	value, err := s.GetAutoDisableErrorThreshold(ctx, projectID)
	if err != nil {
		return 0, err
	}

	s.cache.Set(cacheKey, value, CacheTTL)

	return value, nil
}

func (s *Service) GetAutoDisableTimeWindowSecCached(ctx context.Context, projectID domain.ProjectID) (int, error) {
	cacheKey := makeProjectSettingCacheKey(projectID, AutoDisableTimeWindowSecKey)

	if cached, found := s.cache.Get(cacheKey); found {
		if value, ok := cached.(int); ok {
			return value, nil
		}
	}

	value, err := s.GetAutoDisableTimeWindowSec(ctx, projectID)
	if err != nil {
		return 0, err
	}

	s.cache.Set(cacheKey, value, CacheTTL)

	return value, nil
}

func (s *Service) GetAutoDisableTimeWindowCached(ctx context.Context, projectID domain.ProjectID) (time.Duration, error) {
	cacheKey := makeProjectSettingCacheKey(projectID, AutoDisableTimeWindowSecKey)

	if cached, found := s.cache.Get(cacheKey); found {
		if secs, ok := cached.(int); ok {
			return time.Duration(secs) * time.Second, nil
		}
	}

	value, err := s.GetAutoDisableTimeWindow(ctx, projectID)
	if err != nil {
		return DefaultTimeWindow, err
	}

	s.cache.Set(cacheKey, int(value.Seconds()), CacheTTL)

	return value, nil
}

func (s *Service) InvalidateCache(projectID domain.ProjectID) {
	keys := []string{
		makeProjectSettingCacheKey(projectID, AutoDisableEnabledKey),
		makeProjectSettingCacheKey(projectID, AutoDisableRequiresApprovalKey),
		makeProjectSettingCacheKey(projectID, AutoDisableErrorThresholdKey),
		makeProjectSettingCacheKey(projectID, AutoDisableTimeWindowSecKey),
	}

	for _, key := range keys {
		s.cache.Delete(key)
	}
}

func makeProjectSettingCacheKey(projectID domain.ProjectID, settingName string) string {
	return string(projectID) + ":" + settingName
}

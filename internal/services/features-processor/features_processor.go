package featuresprocessor

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
)

const (
	defaultPollInterval = time.Second * 5
)

type FeaturePrepared struct {
	domain.FeatureExtended

	crons CronsMap
}

type CronsMap map[domain.FeatureScheduleID]cron.Schedule

type ProjectFeatures map[string]FeaturePrepared // key of map is feature.key

type Holder map[domain.ProjectID]ProjectFeatures

type Service struct {
	holder Holder
	mu     sync.RWMutex

	featuresUC   contract.FeaturesUseCase
	projectsUC   contract.ProjectsUseCase
	auditRepo    contract.AuditLogRepository
	pollInterval time.Duration
	lastSeen     time.Time
}

func New(
	featuresUC contract.FeaturesUseCase,
	projectsUC contract.ProjectsUseCase,
	auditRepo contract.AuditLogRepository,
	pollInterval time.Duration,
) *Service {
	if pollInterval <= 0 {
		pollInterval = defaultPollInterval
	}

	return &Service{
		holder:       Holder{},
		featuresUC:   featuresUC,
		projectsUC:   projectsUC,
		auditRepo:    auditRepo,
		pollInterval: pollInterval,
	}
}

func (s *Service) LoadAllFeatures(ctx context.Context) error {
	if s.featuresUC == nil || s.projectsUC == nil {
		return fmt.Errorf("features processor: dependencies not set")
	}

	lastSeen := time.Now()

	projects, err := s.projectsUC.List(ctx)
	if err != nil {
		return fmt.Errorf("list projects: %w", err)
	}

	newHolder := make(Holder, len(projects))
	for _, project := range projects {
		items, err := s.featuresUC.ListExtendedByProjectID(ctx, project.ID)
		if err != nil {
			return fmt.Errorf("list features for project %s: %w", project.ID, err)
		}

		features := make(ProjectFeatures, len(items))
		for _, it := range items {
			features[it.Key] = MakeFeaturePrepared(it)
		}

		newHolder[project.ID] = features
	}

	s.mu.Lock()
	s.holder = newHolder
	s.lastSeen = lastSeen
	s.mu.Unlock()

	return nil
}

func (s *Service) Watch(ctx context.Context) error {
	if s.auditRepo == nil || s.featuresUC == nil {
		return fmt.Errorf("features processor: dependencies not set")
	}

	if s.pollInterval <= 0 {
		s.pollInterval = defaultPollInterval
	}

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	last := s.lastSeen

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			windowEnd := time.Now().UTC()

			logs, err := s.auditRepo.ListSince(ctx, last)
			if err != nil {
				slog.Error("audit log list since failed", "err", err)
				last = windowEnd

				continue
			}

			// Deduplicate by project+feature; keep delete if any delete for feature entity appears.
			type changeKey struct {
				projectID domain.ProjectID
				featureID domain.FeatureID
			}

			changes := make(map[changeKey]domain.AuditAction)
			for _, row := range logs {
				if row.FeatureID == "" {
					continue
				}

				key := changeKey{projectID: row.ProjectID, featureID: row.FeatureID}
				if row.Entity == domain.EntityFeature && row.Action == domain.AuditActionDelete {
					changes[key] = domain.AuditActionDelete

					continue
				}

				if _, ok := changes[key]; !ok {
					changes[key] = domain.AuditActionUpdate
				}
			}

			for key, action := range changes {
				if action == domain.AuditActionDelete {
					ctxRm, cancel := context.WithTimeout(ctx, time.Second*5)
					defer cancel() // TODO: refactor

					s.removeFeatureFromHolder(ctxRm, key.projectID, key.featureID)

					continue
				}

				// refresh feature by loading from repo
				if err := s.refreshFeature(ctx, key.projectID, key.featureID); err != nil {
					slog.Error("refresh feature failed",
						"project", key.projectID, "feature", key.featureID, "err", err)
				}
			}

			last = windowEnd
		}
	}
}

func (s *Service) refreshFeature(ctx context.Context, projectID domain.ProjectID, featureID domain.FeatureID) error {
	featureExtended, err := s.featuresUC.GetExtendedByID(ctx, featureID)
	if err != nil {
		return fmt.Errorf("get feature extended by id %s: %w", featureID, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	featuresMap, ok := s.holder[projectID]
	if !ok {
		featuresMap = ProjectFeatures{}
		s.holder[projectID] = featuresMap
	}

	featuresMap[featureExtended.Key] = MakeFeaturePrepared(featureExtended)

	return nil
}

func (s *Service) removeFeatureFromHolder(
	ctx context.Context,
	projectID domain.ProjectID,
	featureID domain.FeatureID,
) {
	feature, err := s.featuresUC.GetByID(ctx, featureID)
	if err != nil {
		slog.Error("get feature by id failed", "err", err)

		// fallback algorithm
		s.mu.Lock()
		defer s.mu.Unlock()

		if featuresMap, ok := s.holder[projectID]; ok {
			for key, feat := range featuresMap {
				if feat.ID == featureID {
					delete(featuresMap, key)
				}
			}
		}

		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if featuresMap, ok := s.holder[projectID]; ok {
		delete(featuresMap, feature.Key)
	}
}

func (s *Service) Evaluate(
	projectID domain.ProjectID,
	featureKey string,
	reqCtx map[domain.RuleAttribute]any,
) (value string, enabled bool, found bool) {
	feature, ok := s.fetchFeature(projectID, featureKey)
	if !ok {
		return "", false, false
	}

	if !IsFeatureActiveNow(feature, time.Now().UTC()) {
		return "", false, true
	}

	switch feature.Kind {
	case domain.FeatureKindBoolean:
		return feature.DefaultVariant, true, true
	case domain.FeatureKindMultivariant:
		for _, rule := range feature.Rules {
			if !EvaluateExpression(rule.Conditions, reqCtx) {
				continue
			}

			switch rule.Action {
			case domain.RuleActionAssign:
				if rule.FlagVariantID != nil {
					if variant, ok := findVariantByID(feature.FlagVariants, *rule.FlagVariantID); ok {
						return variant.Name, true, true
					}
				} else {
					slog.Error("nil flag variant ID", "rule", rule.ID)
				}
			case domain.RuleActionInclude:
				value = rolloutOrDefault(feature.FlagVariants, feature.RolloutKey, reqCtx, feature.DefaultVariant)

				return value, true, true
			case domain.RuleActionExclude:
				return feature.DefaultVariant, true, true
			}
		}

		value = rolloutOrDefault(feature.FlagVariants, feature.RolloutKey, reqCtx, feature.DefaultVariant)

		return value, true, true
	default:
		return feature.DefaultVariant, true, true
	}
}

func (s *Service) fetchFeature(projectID domain.ProjectID, featureKey string) (FeaturePrepared, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	features, ok := s.holder[projectID]
	if !ok {
		return FeaturePrepared{}, false
	}

	feature, ok := features[featureKey]

	return feature, ok
}

func EvaluateExpression(expr domain.BooleanExpression, reqCtx map[domain.RuleAttribute]any) bool {
	if expr.Condition != nil {
		return MatchCondition(reqCtx, *expr.Condition)
	}

	if expr.Group != nil {
		switch expr.Group.Operator {
		case domain.LogicalOpAND:
			for _, child := range expr.Group.Children {
				if !EvaluateExpression(child, reqCtx) {
					return false
				}
			}

			return true
		case domain.LogicalOpOR:
			for _, child := range expr.Group.Children {
				if EvaluateExpression(child, reqCtx) {
					return true
				}
			}

			return false
		case domain.LogicalOpANDNot:
			if len(expr.Group.Children) != 2 {
				return false
			}

			left := EvaluateExpression(expr.Group.Children[0], reqCtx)
			right := EvaluateExpression(expr.Group.Children[1], reqCtx)

			return left && !right
		}
	}

	return false
}

func IsFeatureActiveNow(feature FeaturePrepared, now time.Time) bool {
	var chosen *domain.FeatureSchedule
	for i := range feature.Schedules {
		schedule := feature.Schedules[i]
		if IsScheduleActive(schedule, feature.crons, now) {
			if chosen == nil {
				chosen = &schedule

				continue
			}
			if schedule.CreatedAt.After(chosen.CreatedAt) {
				chosen = &schedule
			} else if schedule.CreatedAt.Equal(chosen.CreatedAt) {
				// disable важнее enable
				if schedule.Action == domain.FeatureScheduleActionDisable &&
					chosen.Action == domain.FeatureScheduleActionEnable {
					chosen = &schedule
				}
			}
		}
	}

	if chosen != nil {
		return chosen.Action == domain.FeatureScheduleActionEnable
	}

	return feature.Enabled
}

func IsScheduleActive(schedule domain.FeatureSchedule, crons CronsMap, now time.Time) bool {
	loc, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		slog.Error("error loading timezone", "timezone", schedule.Timezone)
		loc = time.UTC
	}
	now = now.In(loc)

	if schedule.StartsAt != nil && now.Before(*schedule.StartsAt) {
		return false
	}
	if schedule.EndsAt != nil && now.After(*schedule.EndsAt) {
		return false
	}

	if schedule.CronExpr == nil || *schedule.CronExpr == "" {
		return true
	}

	sched, ok := crons[schedule.ID]
	if !ok {
		slog.Error("error parsing cron expression", "cron expression", *schedule.CronExpr)

		return false
	}

	prev := sched.Next(now.Add(-time.Minute))
	next := sched.Next(prev)

	return !now.Before(prev) && now.Before(next)
}

func MatchCondition(reqCtx map[domain.RuleAttribute]any, condition domain.Condition) bool {
	actual, ok := reqCtx[condition.Attribute]
	if !ok {
		return false
	}

	switch condition.Operator {
	case domain.OpEq:
		return fmt.Sprint(actual) == fmt.Sprint(condition.Value)
	case domain.OpNotEq:
		return fmt.Sprint(actual) != fmt.Sprint(condition.Value)
	case domain.OpIn:
		return InList(actual, condition.Value, true)
	case domain.OpNotIn:
		return !InList(actual, condition.Value, true)
	case domain.OpGt, domain.OpGte, domain.OpLt, domain.OpLte:
		return CompareNumbers(actual, condition.Value, condition.Operator)
	case domain.OpRegex:
		pattern := fmt.Sprint(condition.Value)
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false
		}

		return re.MatchString(fmt.Sprint(actual))
	case domain.OpPercentage:
		percent, ok := ToInt(condition.Value)
		if !ok {
			return false
		}

		key := fmt.Sprint(actual)
		hash := StableHash(key) % 100

		return hash < percent
	}

	return false
}

func InList(actual any, value any, caseInsensitive bool) bool {
	items, ok := value.([]any)
	if !ok {
		switch v := value.(type) {
		case []string:
			for _, it := range v {
				if caseInsensitive {
					if strings.EqualFold(fmt.Sprint(actual), it) {
						return true
					}
				} else if fmt.Sprint(actual) == it {
					return true
				}
			}
			return false
		default:
			return false
		}
	}

	for _, it := range items {
		if caseInsensitive {
			if strings.EqualFold(fmt.Sprint(actual), fmt.Sprint(it)) {
				return true
			}
		} else if fmt.Sprint(actual) == fmt.Sprint(it) {
			return true
		}
	}

	return false
}

func CompareNumbers(actual any, expected any, op domain.RuleOperator) bool {
	av, aok := ToFloat(actual)
	ev, eok := ToFloat(expected)
	if !aok || !eok {
		return false
	}

	switch op {
	case domain.OpGt:
		return av > ev
	case domain.OpGte:
		return av >= ev
	case domain.OpLt:
		return av < ev
	case domain.OpLte:
		return av <= ev
	}

	return false
}

func ToFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case string:
		f, err := strconv.ParseFloat(n, 64)

		return f, err == nil
	default:
		return 0, false
	}
}

func ToInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	case string:
		i, err := strconv.ParseInt(n, 10, 64)

		return int(i), err == nil
	default:
		return 0, false
	}
}

func StableHash(str string) int {
	hash := 0
	for _, char := range str {
		hash = int(char) + ((hash << 5) - hash)
	}

	if hash < 0 {
		hash = -hash
	}

	return hash
}

func PickVariant(variants []domain.FlagVariant, key string, defaultVariant string) string {
	hash := StableHash(key) % 100
	acc := 0
	for _, v := range variants {
		acc += int(v.RolloutPercent)
		if hash < acc {
			return v.Name
		}
	}

	return defaultVariant
}

func findVariantByID(variants []domain.FlagVariant, id domain.FlagVariantID) (domain.FlagVariant, bool) {
	for _, variant := range variants {
		if variant.ID == id {
			return variant, true
		}
	}

	return domain.FlagVariant{}, false
}

func rolloutOrDefault(
	variants []domain.FlagVariant,
	rolloutKey domain.RuleAttribute,
	reqCtx map[domain.RuleAttribute]any,
	defaultVariant string,
) string {
	if rolloutValue, ok := reqCtx[rolloutKey]; ok {
		return PickVariant(variants, fmt.Sprint(rolloutValue), defaultVariant)
	}

	return defaultVariant
}

func MakeFeaturePrepared(feature domain.FeatureExtended) FeaturePrepared {
	result := FeaturePrepared{
		FeatureExtended: feature,
		crons:           CronsMap{},
	}

	for _, sched := range feature.Schedules {
		if sched.CronExpr != nil && *sched.CronExpr != "" {
			cronSched, err := ParseSchedule(*sched.CronExpr)
			if err != nil {
				slog.Error("invalid cron expr", "expr", *sched.CronExpr, "err", err)

				continue
			}

			result.crons[sched.ID] = cronSched
		}
	}

	return result
}

func ParseSchedule(expr string) (cron.Schedule, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	return parser.Parse(expr)
}

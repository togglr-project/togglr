package featuresprocessor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/rom8726/etoggle/internal/domain"
)

type ProjectFeatures map[string]domain.FeatureExtended

type Holder map[domain.ProjectID]ProjectFeatures

type Service struct {
	holder atomic.Pointer[Holder]
}

func New() *Service {
	return &Service{
		holder: atomic.Pointer[Holder]{},
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

	if !feature.Enabled {
		return "", false, true
	}

	switch feature.Kind {
	case domain.FeatureKindBoolean:
		return feature.DefaultVariant, true, true
	case domain.FeatureKindMultivariant:
		for _, rule := range feature.Rules {
			matched := true
			for _, condition := range rule.Conditions {
				if !MatchCondition(reqCtx, condition) {
					matched = false

					break
				}
			}
			if !matched {
				continue
			}

			switch rule.Action {
			case domain.RuleActionAssign:
				if rule.FlagVariantID != nil {
					if variant, ok := findVariantByID(feature.FlagVariants, *rule.FlagVariantID); ok {
						return variant.Name, true, true
					}
				}
			case domain.RuleActionInclude:
				if userKey, ok := reqCtx[feature.RolloutKey]; ok {
					variant := PickVariant(feature.FlagVariants, fmt.Sprint(userKey), feature.DefaultVariant)

					return variant, true, true
				}

				return feature.DefaultVariant, true, true
			case domain.RuleActionExclude:
				return feature.DefaultVariant, true, true
			}
		}

		if userKey, ok := reqCtx[feature.RolloutKey]; ok {
			variant := PickVariant(feature.FlagVariants, fmt.Sprint(userKey), feature.DefaultVariant)

			return variant, true, true
		}

		return feature.DefaultVariant, true, true
	default:
		return feature.DefaultVariant, true, true
	}
}

func (s *Service) fetchFeature(projectID domain.ProjectID, featureKey string) (domain.FeatureExtended, bool) {
	holder := s.holder.Load()
	features, ok := (*holder)[projectID]
	if !ok {
		return domain.FeatureExtended{}, false
	}

	feature, ok := features[featureKey]

	return feature, ok
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

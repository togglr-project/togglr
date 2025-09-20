package featuresprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/etoggle/internal/domain"
)

func TestService_Evaluate_Table2(t *testing.T) {
	variantA := domain.FlagVariant{ID: domain.FlagVariantID("v1"), Name: "A", RolloutPercent: 50}
	variantB := domain.FlagVariant{ID: domain.FlagVariantID("v2"), Name: "B", RolloutPercent: 50}

	tests := []struct {
		name      string
		fe        domain.FeatureExtended
		reqCtx    map[domain.RuleAttribute]any
		wantValue string // if an empty string expected value is ""
		wantEn    bool
		wantFound bool
		// optional allowed values set for rollout (if wantValueEmpty==false and multiple possibilities)
		allowedValues []string
	}{
		{
			name: "exclude rule disables feature",
			fe: makeFeatureExtended([]domain.Rule{
				{
					ID:     "r1",
					Action: domain.RuleActionExclude,
					Conditions: domain.BooleanExpression{
						Condition: &domain.Condition{
							Attribute: "country",
							Operator:  domain.OpEq,
							Value:     "RU",
						},
					},
					Priority: 0,
				},
			}, nil, domain.FeatureKindSimple, "default"),
			reqCtx:        map[domain.RuleAttribute]any{"country": "RU"},
			wantValue:     "",
			wantEn:        false,
			wantFound:     true,
			allowedValues: nil,
		},
		{
			name: "assign rule selects variant",
			fe: makeFeatureExtended([]domain.Rule{
				{
					ID:            "r1",
					Action:        domain.RuleActionAssign,
					FlagVariantID: ptrFV("v1"),
					Conditions: domain.BooleanExpression{
						Condition: &domain.Condition{
							Attribute: "country",
							Operator:  domain.OpEq,
							Value:     "RU",
						},
					},
					Priority: 0,
				},
			}, []domain.FlagVariant{variantA}, domain.FeatureKindMultivariant, "default"),
			reqCtx:        map[domain.RuleAttribute]any{"country": "RU"},
			wantValue:     "A",
			wantEn:        true,
			wantFound:     true,
			allowedValues: nil,
		},
		{
			name: "priority: higher priority assign wins (lower numeric = higher priority)",
			fe: makeFeatureExtended([]domain.Rule{
				{
					ID:            "low",
					Action:        domain.RuleActionAssign,
					FlagVariantID: ptrFV("v1"),
					Priority:      10,
					Conditions: domain.BooleanExpression{
						Condition: &domain.Condition{
							Attribute: "country",
							Operator:  domain.OpEq,
							Value:     "RU",
						},
					},
				},
				{
					ID:            "high",
					Action:        domain.RuleActionAssign,
					FlagVariantID: ptrFV("v2"),
					Priority:      1, // higher priority (smaller number)
					Conditions: domain.BooleanExpression{
						Condition: &domain.Condition{
							Attribute: "country",
							Operator:  domain.OpEq,
							Value:     "RU",
						},
					},
				},
			}, []domain.FlagVariant{variantA, variantB}, domain.FeatureKindMultivariant, "default"),
			reqCtx:    map[domain.RuleAttribute]any{"country": "RU"},
			wantValue: "B",
			wantEn:    true,
			wantFound: true,
		},
	}

	svc := New(nil, nil, nil, 0)
	// initialize holder map
	svc.mu.Lock()
	svc.holder = Holder{}
	svc.mu.Unlock()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare FeaturePrepared using MakeFeaturePrepared
			fp := MakeFeaturePrepared(tt.fe)

			// insert into holder
			svc.mu.Lock()
			if svc.holder == nil {
				svc.holder = Holder{}
			}
			svc.holder[tt.fe.ProjectID] = ProjectFeatures{
				tt.fe.Key: fp,
			}
			svc.mu.Unlock()

			value, enabled, found := svc.Evaluate(tt.fe.ProjectID, tt.fe.Key, tt.reqCtx)
			assert.Equal(t, tt.wantEn, enabled, "enabled mismatch")
			assert.Equal(t, tt.wantFound, found, "found mismatch")

			if tt.allowedValues != nil {
				// value must be one of allowedValues
				assert.Contains(t, tt.allowedValues, value)
			} else {
				assert.Equal(t, tt.wantValue, value)
			}
		})
	}
}

func ptrFV(id string) *domain.FlagVariantID {
	v := domain.FlagVariantID(id)
	return &v
}

func makeFeatureExtended(rules []domain.Rule, variants []domain.FlagVariant, kind domain.FeatureKind, defaultVariant string) domain.FeatureExtended {
	return domain.FeatureExtended{
		Feature: domain.Feature{
			ID:             "f1",
			ProjectID:      "p1",
			Key:            "feature_key",
			Name:           "Test Feature",
			Kind:           kind,
			DefaultVariant: defaultVariant,
			Enabled:        true,
			CreatedAt:      time.Now(),
		},
		FlagVariants: variants,
		Rules:        rules,
		// no schedules in tests (nil)
	}
}

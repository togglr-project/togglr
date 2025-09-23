//go:build integration

package tests

import (
	"testing"

	"github.com/togglr-project/togglr/tests/runner"
)

func TestFeatureSchedulesAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/feature_schedules",
	}
	runner.Run(t, &cfg)
}

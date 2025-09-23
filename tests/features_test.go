//go:build integration

package tests

import (
	"testing"

	"github.com/togglr-project/togglr/tests/runner"
)

func TestFeaturesAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/features",
	}
	runner.Run(t, &cfg)
}

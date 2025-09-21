//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/etoggle/tests/runner"
)

func TestFeaturesAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/features",
	}
	runner.Run(t, &cfg)
}

//go:build integration

package tests

import (
	"testing"

	"github.com/togglr-project/togglr/tests/runner"
)

func TestSDKReportErrorAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/sdk/report-error",
	}
	runner.RunSDK(t, &cfg)
}

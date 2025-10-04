//go:build integration

package tests

import (
	"testing"
	"time"

	"github.com/togglr-project/togglr/internal"
	"github.com/togglr-project/togglr/tests/runner"
)

func TestSDKReportErrorAPI(t *testing.T) {
	t.SkipNow()
	cfg := runner.Config{
		CasesDir: "./cases/sdk/report-error",
		AfterReq: func(app *internal.App) error {
			time.Sleep(2 * time.Second)

			return nil
		},
	}
	runner.RunSDK(t, &cfg)
}

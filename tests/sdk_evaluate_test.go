//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/etoggle/tests/runner"
)

func TestSDKEvaluateAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/sdk/evaluate",
	}
	runner.RunSDK(t, &cfg)
}

//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/etoggle/tests/runner"
)

func TestRuleAttributesAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/rule_attributes",
	}
	runner.Run(t, &cfg)
}

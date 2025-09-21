//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/etoggle/tests/runner"
)

func TestProjectsAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/projects",
	}
	runner.Run(t, &cfg)
}

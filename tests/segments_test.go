//go:build integration

package tests

import (
	"testing"

	"github.com/togglr-project/togglr/tests/runner"
)

func TestSegmentsAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/segments",
	}
	runner.Run(t, &cfg)
}

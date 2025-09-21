//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/etoggle/tests/runner"
)

func TestSegmentsAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/segments",
	}
	runner.Run(t, &cfg)
}

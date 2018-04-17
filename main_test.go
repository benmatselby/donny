package main

import (
	"os"
	"testing"
)

func TestEnvironmentVars(t *testing.T) {
	os.Unsetenv("VSTS_TOKEN")
	os.Unsetenv("VSTS_ACCOUNT")
	os.Unsetenv("VSTS_PROJECT")
	os.Unsetenv("VSTS_TEAM")

	_, err := loadEnvironmentVars()
	if err == nil {
		t.Errorf("expected error")
	}
}

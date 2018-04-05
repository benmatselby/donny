package vsts

import (
	"testing"
)

func TestClient_New(t *testing.T) {
	c := New(
		"my-account",
		"my-project",
		"my-team",
		"my-token",
	)

	if c.Account != "my-account" {
		t.Errorf("Client.Account = %s; expected %s", c.Account, "my-account")
	}
}

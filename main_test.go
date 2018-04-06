package main

import (
	"fmt"
	"os"
	"testing"
)

func TestEnvironmentVars(t *testing.T) {
	os.Unsetenv("VSTS_TOKEN")
	os.Unsetenv("VSTS_ACCOUNT")
	os.Unsetenv("VSTS_PROJECT")
	os.Unsetenv("VSTS_TEAM")

	expected := `
In order for donny to integrate with VSTS, you need to define the following environment variables:

* VSTS_ACCOUNT = %s
* VSTS_PROJECT = %s
* VSTS_TEAM    = %s
* VSTS_TOKEN   = %s
`
	_, err := environmentVars()
	expectedErr := fmt.Errorf(expected, "", "", "", "")

	if expectedErr.Error() != err.Error() {
		t.Errorf("\nmessage = %s\nexpected = %s", err, expectedErr)
	}
}

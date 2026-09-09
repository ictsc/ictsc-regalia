package config

import "testing"

func TestLoadRuntimeRejectsInsecureProductionCookies(t *testing.T) {
	setValidProductionRuntimeEnv(t)
	t.Setenv("ICTSC_SECURE_COOKIES", "false")
	if _, err := LoadRuntime(); err == nil {
		t.Fatal("expected production Secure cookie validation error")
	}
}

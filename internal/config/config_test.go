package config

import "testing"

func TestVersionUsesRuntimeEnvWhenBuildVersionIsDefault(t *testing.T) {
	originalVersion := version
	version = "dev"
	defer func() {
		version = originalVersion
	}()

	t.Setenv("APP_VERSION", "1.2.0")

	if got := Version(); got != "1.2.0" {
		t.Fatalf("expected runtime APP_VERSION to be used, got %q", got)
	}
}

func TestVersionPrefersEmbeddedBuildVersion(t *testing.T) {
	originalVersion := version
	version = "2.3.4"
	defer func() {
		version = originalVersion
	}()

	t.Setenv("APP_VERSION", "1.2.0")

	if got := Version(); got != "2.3.4" {
		t.Fatalf("expected embedded build version to win, got %q", got)
	}
}

func TestVersionFallsBackToDevWhenUnset(t *testing.T) {
	originalVersion := version
	version = ""
	defer func() {
		version = originalVersion
	}()

	t.Setenv("APP_VERSION", "")

	if got := Version(); got != "dev" {
		t.Fatalf("expected dev fallback, got %q", got)
	}
}

package executor

import (
	"context"
	"errors"
	"testing"
)

func TestResolveEnvironmentVariables_ContinueOnError_PartialResults(t *testing.T) {
	envVars := []string{
		"NORMAL=ok",
		"GOOD=@Microsoft.KeyVault(VaultName=v;SecretName=good)",
		"BAD=@Microsoft.KeyVault(VaultName=v;SecretName=missing)",
		"AFTER=akvs://guid/v/after",
		"TAIL=done",
	}

	resolve := func(_ context.Context, reference string) (string, error) {
		switch reference {
		case "@Microsoft.KeyVault(VaultName=v;SecretName=good)":
			return "good-value", nil
		case "akvs://guid/v/after":
			return "after-value", nil
		default:
			return "", errors.New("not found")
		}
	}

	resolved, warnings, err := resolveEnvironmentVariables(context.Background(), envVars, resolve, ResolveEnvironmentOptions{StopOnError: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Key != "BAD" {
		t.Fatalf("expected warning for key BAD, got %q", warnings[0].Key)
	}

	want := []string{
		"NORMAL=ok",
		"GOOD=good-value",
		"BAD=@Microsoft.KeyVault(VaultName=v;SecretName=missing)",
		"AFTER=after-value",
		"TAIL=done",
	}

	if len(resolved) != len(want) {
		t.Fatalf("expected %d resolved vars, got %d", len(want), len(resolved))
	}
	for i := range want {
		if resolved[i] != want[i] {
			t.Fatalf("resolved[%d] = %q, want %q", i, resolved[i], want[i])
		}
	}
}

func TestResolveEnvironmentVariables_StopOnError_FailFast_NoPartialResults(t *testing.T) {
	envVars := []string{
		"GOOD=@Microsoft.KeyVault(VaultName=v;SecretName=good)",
		"BAD=@Microsoft.KeyVault(VaultName=v;SecretName=missing)",
		"AFTER=@Microsoft.KeyVault(VaultName=v;SecretName=after)",
	}

	resolve := func(_ context.Context, reference string) (string, error) {
		if reference == "@Microsoft.KeyVault(VaultName=v;SecretName=good)" {
			return "good-value", nil
		}
		return "", errors.New("not found")
	}

	resolved, warnings, err := resolveEnvironmentVariables(context.Background(), envVars, resolve, ResolveEnvironmentOptions{StopOnError: true})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if resolved != nil {
		t.Fatalf("expected no partial results, got %v", resolved)
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Key != "BAD" {
		t.Fatalf("expected warning for key BAD, got %q", warnings[0].Key)
	}
}

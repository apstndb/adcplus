package adcplus

import (
	"reflect"
	"testing"

	"github.com/apstndb/adcplus/internal/config"
)

func TestWithTargetPrincipal(t *testing.T) {
	var cfg config.AdcPlusConfig
	opt := WithTargetPrincipal("sa@p.iam.gserviceaccount.com")
	if err := opt(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TargetPrincipal != "sa@p.iam.gserviceaccount.com" {
		t.Errorf("TargetPrincipal = %q, want %q", cfg.TargetPrincipal, "sa@p.iam.gserviceaccount.com")
	}
}

func TestWithDelegates(t *testing.T) {
	var cfg config.AdcPlusConfig
	opt := WithDelegates("d1@p.iam.gserviceaccount.com", "d2@p.iam.gserviceaccount.com")
	if err := opt(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"d1@p.iam.gserviceaccount.com", "d2@p.iam.gserviceaccount.com"}
	if !reflect.DeepEqual(cfg.Delegates, want) {
		t.Errorf("Delegates = %v, want %v", cfg.Delegates, want)
	}
}

func TestWithScopes(t *testing.T) {
	var cfg config.AdcPlusConfig
	opt := WithScopes("scope1", "scope2")
	if err := opt(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"scope1", "scope2"}
	if !reflect.DeepEqual(cfg.Scopes, want) {
		t.Errorf("Scopes = %v, want %v", cfg.Scopes, want)
	}
}

func TestWithCredentialsFile(t *testing.T) {
	var cfg config.AdcPlusConfig
	opt := WithCredentialsFile("/path/to/creds.json")
	if err := opt(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CredentialsFile != "/path/to/creds.json" {
		t.Errorf("CredentialsFile = %q, want %q", cfg.CredentialsFile, "/path/to/creds.json")
	}
}

func TestWithCredentialsJSON(t *testing.T) {
	var cfg config.AdcPlusConfig
	j := []byte(`{"type":"service_account"}`)
	opt := WithCredentialsJSON(j)
	if err := opt(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cfg.CredentialsJSON, j) {
		t.Errorf("CredentialsJSON = %s, want %s", cfg.CredentialsJSON, j)
	}
}

func TestWithExperimentalAppEngineSigner(t *testing.T) {
	var cfg config.AdcPlusConfig
	opt := WithExperimentalAppEngineSigner(true)
	if err := opt(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.EnableAppEngineSigner {
		t.Error("EnableAppEngineSigner = false, want true")
	}
}

package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/apstndb/adcplus/internal/config"
	"github.com/apstndb/adcplus/option"
	"golang.org/x/oauth2/google"
)

const impSaEnvName = "CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT"

// All patterns are defined in
// https://github.com/golang/oauth2/blob/f6687ab2804cbebdfdeef385bee94918b1ce83de/google/google.go#L93-L98
const (
	ServiceAccountKey         = "service_account"
	UserCredentialsKey        = "authorized_user"
	ExternalAccountKey        = "external_account"
	ComputeMetadataCredential = "compute_metadata"
)

func CalcSmartSignerConfig(opts ...option.Option) (*config.AdcPlusConfig, error) {
	var config config.AdcPlusConfig
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	if config.TargetPrincipal == "" && len(config.Delegates) > 0 {
		return nil, fmt.Errorf("targetPrincipal is set but delegates is not set: %s", config.Delegates)
	}
	if impSaVal := os.Getenv(impSaEnvName); config.TargetPrincipal == "" && impSaVal != "" {
		config.TargetPrincipal, config.Delegates = ParseDelegateChain(impSaVal)
	}
	return &config, nil
}

func CredentialTypeFromJSON(credentialJSON []byte) (string, error) {
	// Minimal subset of google.credentialsFile
	// https://github.com/golang/oauth2/blob/f6687ab2804cbebdfdeef385bee94918b1ce83de/google/google.go#L100-L126
	type credentialFile struct {
		Type string `json:"type"`
	}
	var parsedCredential credentialFile
	err := json.Unmarshal(credentialJSON, &parsedCredential)
	if err != nil {
		return "", err
	}
	if parsedCredential.Type == "" {
		return "", errors.New("credential type is empty")
	}

	return parsedCredential.Type, nil
}

func InferADCCredentialType(cred *google.Credentials) (string, error) {
	if cred.JSON != nil {
		return CredentialTypeFromJSON(cred.JSON)
	}
	if metadata.OnGCE() {
		return ComputeMetadataCredential, nil
	}
	return "", errors.New("unknown credential type")
}

// ParseDelegateChain split impersonate target principal and delegate chain.
// s must be non-empty string.
func ParseDelegateChain(s string) (targetPrincipal string, delegates []string) {
	if s == "" {
		panic("parseDelegateChain: empty argument")
	}
	ss := strings.Split(s, ",")
	return ss[len(ss)-1], ss[:len(ss)-1]
}

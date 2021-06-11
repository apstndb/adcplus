package signer

import (
	"encoding/json"
	"errors"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2/google"
)

func credentialTypeFromJSON(credentialJSON []byte) (string, error) {
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

func inferADCCredentialType(cred *google.Credentials) (string, error) {
	if cred.JSON != nil {
		return credentialTypeFromJSON(cred.JSON)
	}
	if metadata.OnGCE() {
		return computeMetadataCredential, nil
	}
	return "", errors.New("unknown credential type")
}

// parseDelegateChain split impersonate target principal and delegate chain.
// s must be non-empty string.
func parseDelegateChain(s string) (targetPrincipal string, delegates []string) {
	if s == "" {
		panic("parseDelegateChain: empty argument")
	}
	ss := strings.Split(s, ",")
	return ss[len(ss)-1], ss[:len(ss)-1]
}

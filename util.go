package signer

import (
	"encoding/json"
	"strings"
)

func credentialType(credentialJSON []byte) (string, error) {
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

	return parsedCredential.Type, nil
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

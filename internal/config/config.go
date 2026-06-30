package config

import "golang.org/x/oauth2"

type AdcPlusConfig struct {
	TargetPrincipal       string
	Delegates             []string
	EnableAppEngineSigner bool
	Scopes                []string
	CredentialsFile       string
	CredentialsJSON       []byte
	TokenSource           oauth2.TokenSource
}

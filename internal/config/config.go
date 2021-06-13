package config

type AdcPlusConfig struct {
	TargetPrincipal       string
	Delegates             []string
	EnableAppEngineSigner bool
	Scopes                []string
	CredentialsFile       string
	CredentialsJSON       []byte
}

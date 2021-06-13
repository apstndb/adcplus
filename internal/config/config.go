package config

type AdcPlusConfig struct {
	TargetPrincipal       string
	Delegates             []string
	EnableAppengineSigner bool
	Scopes                []string
	CredentialsFile       string
	CredentialsJSON       []byte
}

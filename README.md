# adcplus

[![Go Reference](https://pkg.go.dev/badge/github.com/apstndb/adcplus.svg)](https://pkg.go.dev/github.com/apstndb/adcplus)

This package implements oauth2.TokenSource and signer which respects [ADC](https://google.aip.dev/auth/4110) with impersonation.
It automatically uses `CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT` environment variable as the impersonation target and the delegation chain.

**This package is EXPERIMENTAL**.

## Underlying method

* "Credentials API" is Service Account Credentials API ([`projects.serviceAccounts.signBlob`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signBlob?hl=en), [`projects.serviceAccounts.signJwt`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signJwt?hl=en))
  * Need [Service Account Token Creator role(`roles/iam.serviceAccountTokenCreator`)](https://cloud.google.com/iam/docs/impersonating-service-accounts)

### signer.SmartSigner

|credential/impersonate|yes|no|
|---|---|---|
|authorized_user|Credentials API|Not Supported|
|service_account|Credentials API|Sign by JSON key|
|external_account|Credentials API|Credentials API as itself|
|compute_metadata|Credentials API|Credentials API as itself|
|App Engine 1st gen(only if `WithExperimentalAppEngineSigner(true)`)|Credentials API|`appengine.SignBytes()`|

### tokensource.SmartAccessTokenSource

|credential/impersonate|yes|no|
|---|---|---|
|authorized_user|Credentials API|ADC(refresh token flow)|
|service_account|Credentials API|ADC(jwt-bearer token flow)|
|external_account|Credentials API|ADC(STS)|
|compute_metadata|Credentials API|ADC(token endpoint)|

### tokensource.SmartIDTokenSource

|credential/impersonate|yes|no|
|---|---|---|
|authorized_user|Credentials API|Not Supported|
|service_account|Credentials API|ADC(jwt-bearer flow)|
|external_account|Credentials API|Not Supported(TODO)|
|compute_metadata|Credentials API|ADC(identity endpoint)|

## TODO

* Support to override underlying TokenSource.
  * `WithTokenSource()`
* Support external_account in `tokensource.SmartIDTokenSource`

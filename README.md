# adcplus

[![Go Reference](https://pkg.go.dev/badge/github.com/apstndb/adcplus.svg)](https://pkg.go.dev/github.com/apstndb/adcplus)

This package implements oauth2.TokenSource and signer which respects [ADC](https://google.aip.dev/auth/4110) with impersonation.

* Automatically uses `CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT` environment variable as an impersonation target and a delegation chain.
  * It respects same variable and syntax of gcloud.
    * https://cloud.google.com/sdk/gcloud/reference/topic/configurations?hl=en#impersonate_service_account
    * https://cloud.google.com/sdk/docs/properties?hl=en#setting_properties_via_environment_variables
* Can override the impersonation target, the delegate chain and the source credential through [functional options](https://pkg.go.dev/github.com/apstndb/adcplus#Option).

## Disclaimer

**This package is EXPERIMENTAL.**
* No responsibility.
* May be broken.
* Will do breaking changes.

## Underlying method

* Currently, external_account(STS) is not mentioned in [AIP-4110](https://google.aip.dev/auth/4110) because it is [removed when approval](https://github.com/aip-dev/google.aip.dev/pull/592) but it is supported in [`golang.org/x/oauth2/google`](https://github.com/golang/oauth2/pull/462) and it is [documented](https://cloud.google.com/docs/authentication/production?hl=en). I treat it as one of ADC credential.
* "Credentials API" is Service Account Credentials API ([`projects.serviceAccounts.signBlob`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signBlob?hl=en), [`projects.serviceAccounts.signJwt`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signJwt?hl=en))
  * Need [Service Account Token Creator role(`roles/iam.serviceAccountTokenCreator`)](https://cloud.google.com/iam/docs/impersonating-service-accounts)

### [signer.SmartSigner](https://pkg.go.dev/github.com/apstndb/adcplus/signer#SmartSigner)

|credential/impersonate|yes|no|
|---|---|---|
|authorized_user|Credentials API|Not Supported|
|service_account|Credentials API|Sign by JSON key|
|external_account|Credentials API|Credentials API as itself|
|compute_metadata|Credentials API|Credentials API as itself|
|App Engine 1st gen(only if `WithExperimentalAppEngineSigner(true)`)|Credentials API|`appengine.SignBytes()`|

### [tokensource.SmartAccessTokenSource](https://pkg.go.dev/github.com/apstndb/adcplus/tokensource#SmartAccessTokenSource)

|credential/impersonate|yes|no|
|---|---|---|
|authorized_user|Credentials API|ADC(refresh token flow)|
|service_account|Credentials API|ADC(jwt-bearer token flow)|
|external_account|Credentials API|ADC(STS)|
|compute_metadata|Credentials API|ADC(token endpoint)|

### [tokensource.SmartIDTokenSource](https://pkg.go.dev/github.com/apstndb/adcplus/tokensource#SmartIDTokenSource)

|credential/impersonate|yes|no|
|---|---|---|
|authorized_user|Credentials API|Not Supported|
|service_account|Credentials API|ADC(jwt-bearer flow)|
|external_account|Credentials API|Not Supported(TODO: retrieve using STS)|
|compute_metadata|Credentials API|ADC(identity endpoint)|

## TODO

* Support to override underlying TokenSource.
  * `WithTokenSource()`
* Support external_account in `tokensource.SmartIDTokenSource`.
* Re-implement underlying TokenSource to avoid ReuseTokenSource in default.
* Add tests.

# adcplus

[![Go Reference](https://pkg.go.dev/badge/github.com/apstndb/adcplus.svg)](https://pkg.go.dev/github.com/apstndb/adcplus)

**This package is EXPERIMENTAL**.

## signer.SmartSigner

### Underlying signing method

|credential/impersonate|yes|no|
|---|---|---|
|service_account|Credentials API|Sign by JSON key|
|authorized_user|Credentials API|Not Supported|
|external_account|Credentials API|Credentials API as itself|
|compute_metadata|Credentials API|Credentials API as itself|
|App Engine 1st gen(only if `WithExperimentalAppEngineSigner(true)`)|Credentials API|`appengine.SignBytes()`|

* "Credentials API" is Service Account Credentials API ([`projects.serviceAccounts.signBlob`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signBlob?hl=en), [`projects.serviceAccounts.signJwt`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signJwt?hl=en))
  * Need [Service Account Token Creator role(`roles/iam.serviceAccountTokenCreator`)](https://cloud.google.com/iam/docs/impersonating-service-accounts)

## tokensource.SmartAccessTokenSource

### Underlying token source

|credential/impersonate|yes|no|
|---|---|---|
|service_account|Credentials API|ADC(jwt-bearer token flow)|
|authorized_user|Credentials API|ADC(refresh token flow)|
|external_account|Credentials API|ADC(STS)|
|compute_metadata|Credentials API|ADC(token endpoint)|

## tokensource.SmartIDTokenSource

### Underlying token source

|credential/impersonate|yes|no|
|---|---|---|
|service_account|Credentials API|ADC(jwt-bearer flow)|
|authorized_user|Credentials API|Not Supported|
|external_account|Credentials API|Not Supported(TODO)|
|compute_metadata|Credentials API|ADC(identity endpoint)|

## TODO

* Support to override underlying TokenSource.
  * `WithCredentialsFile()`
  * `WithCredentialsJSON()`
  * `WithTokenSource()`
* Support external_account in `tokensource.SmartIDTokenSource`

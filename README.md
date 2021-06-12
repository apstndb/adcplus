# signer

[![Go Reference](https://pkg.go.dev/badge/github.com/apstndb/signer.svg)](https://pkg.go.dev/github.com/apstndb/signer)

**This package is EXPERIMENTAL**.

## Underlying signing method

|credential/impersonate|yes|no|
|---|---|---|
|service_account|Credentials API|Sign by JSON key|
|authorized_user|Credentials API|Not Supported|
|external_account|Credentials API|Credentials API as itself|
|compute_metadata|Credentials API|Credentials API as itself|
|App Engine 1st gen(only if `WithExperimentalAppEngineSigner(true)`)|Credentials API|`appengine.SignBytes()`|

* "Credentials API" is Service Account Credentials API ([`projects.serviceAccounts.signBlob`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signBlob?hl=en), [`projects.serviceAccounts.signJwt`](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signJwt?hl=en))
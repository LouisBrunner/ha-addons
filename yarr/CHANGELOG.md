# Changelog

## 0.0.3

- Validate `hostname`, `oidc.issuer`, and `base_path` against their expected formats
- Stop logging the full add-on configuration on startup
- `fever_password` is now an optional config option instead of an auto-generated credential you had to dig out of `/data/fever_auth`; leave it unset to disable Fever API access entirely

## 0.0.2

- Fever API support for mobile clients, bypassing OIDC via a separate generated credential

## 0.0.1

Initial release

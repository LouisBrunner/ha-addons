# Zot

Self-hosted [Zot](https://zotregistry.dev) OCI/Docker container registry, based on the official `ghcr.io/project-zot/zot` image which bundles the [ZUI](https://github.com/project-zot/zui) web interface.

## Prerequisites

- A hostname pointing at your Home Assistant for HTTPS routing (Cloudflare tunnel, Nginx reverse proxy, etc).
- An OIDC provider with a client created for this instance (e.g. [Pocket ID](https://pocket-id.org)), callback URL `<external_url>/zot/auth/callback/oidc`.

## Configuration

| Option                           | Required | Description                                                                              |
| -------------------------------- | -------- | ---------------------------------------------------------------------------------------- |
| `external_url`                   | Yes      | Public URL of this Zot instance, e.g. `https://zot.mydomain.com`                         |
| `blob_upload_timeout`            | No       | Max time per blob chunk read/write before aborting (default: `2m`)                       |
| `oidc.issuer`                    | Yes      | Your OIDC provider's issuer URL                                                          |
| `oidc.client_id`                 | Yes      | OIDC client ID                                                                           |
| `oidc.client_secret`             | Yes      | OIDC client secret                                                                       |
| `oidc.provider_name`             | No       | Display name on the login button (default: empty)                                        |
| `oidc.username_claim`            | No       | Claim used as identity for access control (default: `email`)                             |
| `oidc.admin_identities`          | No       | Values of the username claim granted full admin access                                   |
| `oidc.admin_groups`              | No       | OIDC group names granted full admin access                                               |
| `anonymous_read`                 | No       | Allow unauthenticated pulls (default: `false`)                                           |
| `readonly_groups`                | No       | OIDC group names (from the `groups` claim) with global read access                       |
| `readonly_identities`            | No       | Values of the username claim with global read access                                     |
| `repositories`                   | No       | List of per-repository access overrides, see [Authorization model](#authorization-model) |
| `storage.gc`                     | No       | Enable garbage collection of unreferenced blobs (default: `true`)                        |
| `storage.gc_delay`               | No       | Delay before an unreferenced blob is collected (default: `1h`)                           |
| `storage.dedupe`                 | No       | Deduplicate identical blobs across repositories (default: `true`)                        |
| `storage.retention`              | No       | Tag/image retention policies, see [Retention](#retention)                                |
| `metrics_enabled`                | No       | Expose a Prometheus `/metrics` endpoint (default: `false`)                               |
| `cve.enabled`                    | No       | Scan pushed images for CVEs via Trivy, needs outbound internet (default: `true`)         |
| `cve.update_interval`            | No       | How often to refresh the vulnerability database (default: `24h`)                         |
| `rate_limit.enabled`             | No       | Enable a global request rate limit (default: `true`)                                     |
| `rate_limit.requests_per_second` | No       | Max requests/second across all methods not overridden below (default: `100`)             |
| `rate_limit.per_method`          | No       | List of `{method, requests_per_second}` overrides for specific HTTP methods              |
| `webhooks`                       | No       | List of HTTP endpoints notified of registry events, see [Webhooks](#webhooks)            |
| `lint.enabled`                   | No       | Reject pushes missing mandatory annotations (default: `false`)                           |
| `lint.mandatory_annotations`     | No       | List of required annotation keys, e.g. `org.opencontainers.image.source`                 |
| `scrub.enabled`                  | No       | Enable periodic blob integrity checks (default: `false`)                                 |
| `scrub.interval`                 | No       | Time between scrub runs, min `2h` (default: `24h`)                                       |
| `sync`                           | No       | List of upstream registries to mirror, see [Sync (mirroring)](#sync-mirroring)           |
| `debug`                          | No       | Enable extra diagnostic logging (default: `false`)                                       |

## Ports

| Port       | Direction | Description                                   |
| ---------- | --------- | --------------------------------------------- |
| `5000/tcp` | inbound   | Registry API + web UI (point your proxy here) |

Not exposed to the host by default, assuming a HTTPS reverse proxy in front (e.g. `http://{SLUG}-zot:5000`).

## Persistent data

| Path                      | Contents                                                    |
| ------------------------- | ----------------------------------------------------------- |
| `/data/registry/*.db`     | Databases (`storage.rootDirectory`)                         |
| `/data/session-keys.json` | OIDC session hash/encrypt keys, generated once on first run |

> ![!IMPORTANT]
> Images and session tokens are not backed up by Home Assistant.

## Authentication

Web UI login is OIDC only. For `docker login`/push/pull, generate an API key from the account menu after logging in (username must match your `oidc.username_claim` value):

```sh
docker login your-hostname -u your-username-claim-value -p zak_theSelfGeneratedKey
```

If the OIDC provider is down, existing keys keep working but no one can log in or mint new ones.

## Authorization model

- Identities and groups come straight from the OIDC provider on every login, via `oidc.username_claim` and the `groups` claim. Nothing to provision locally.
- `oidc.admin_identities`/`admin_groups` get full CRUD on all repositories regardless of anything below (Zot's `adminPolicy`).
- `readonly_identities`/`readonly_groups` and `anonymous_read` are the global/default policy, applied to any repository not matched by a `repositories` entry.
- `repositories` overrides per pattern:
  ```yaml
  repositories:
    - pattern: team/**
      readonly_identities:
        - ci
      readonly_groups:
        - viewers
      anonymous_read: false
    - pattern: public/**
      readonly_identities: []
      readonly_groups: []
      anonymous_read: true
  ```
  `readonly_identities` are username-claim values, `readonly_groups` are OIDC group names, both matched live against the session.

> [!IMPORTANT]
> A matching `repositories` pattern fully replaces the global (`**`) policy for that repo, it doesn't add to it. A user only in the global `readonly_identities` gets `403` on a repo matched by a pattern that doesn't also list them. List them in both places if you want both.

## Webhooks

Each entry is an `http` event sink, Zot posts CloudEvents payloads to `url` for pushes/deletions/etc:

```yaml
webhooks:
  - url: https://example.com/hook
    username: someuser
    password: somepassword
  - url: https://example.com/hook2
    token: sometoken
```

`username`+`password` or `token`, not both.

## Sync (mirroring)

Each entry mirrors one upstream. `on_demand: true` pulls on first request and caches. Add `poll_interval`/`content_prefix` to also mirror a prefix periodically:

```yaml
sync:
  - url: https://registry-1.docker.io
    on_demand: true
  - url: https://ghcr.io
    on_demand: false
    poll_interval: 6h
    content_prefix: myorg/**
    content_destination: mirrored/myorg
```

- TLS is always verified against upstreams, no option to disable it.
- One URL per upstream and one prefix/destination pair per entry here.
- `username`/`password` authenticate to that upstream registry.

## Retention

Different from `storage.gc`: retention deletes tags/images still referenced, by rule. Each `keep_tags` entry is evaluated per matching repository, a tag survives if it matches any entry:

```yaml
storage:
  retention:
    enabled: true
    delay: 24h
    policies:
      - repositories:
          - staging/**
        delete_untagged: true
        delete_referrers: true
        keep_tags:
          - patterns:
              - v.*
            most_recently_pushed_count: 10
          - patterns:
              - latest
```

- Deletes for real, no dry-run mode, test your patterns against a non-critical repository first.
- `delay` is separate from `storage.gc_delay`, it's how long to wait before removing untagged images/referrers.
- `keep_tags[].patterns` are regexes, not globs.
- `pulled_within`/`pushed_within` (durations) and `most_recently_pushed_count` further restrict a `keep_tags` entry.
- Repos not matched by any policy are unaffected.

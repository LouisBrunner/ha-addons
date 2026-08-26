# yarr

Self-hosted [yarr](https://github.com/nkanaev/yarr) ("yet another RSS reader"), a minimal RSS/Atom feed aggregator:
single Go binary, embedded SQLite, no external database, no themes/extensions/feature creep.

Login is via OIDC only, set your provider's details below and there's nothing else to run or configure.

## Prerequisites

- A public hostname pointing at your Home Assistant (e.g. via Cloudflare tunnel, Nginx reverse proxy, etc.) with HTTPS in front.
- An OIDC provider with a client created for this instance (e.g. [Pocket ID](https://pocket-id.org)), redirect URI `https://<hostname>/oauth2/callback`.

## Configuration

| Option               | Required | Description                                                                                                     |
| -------------------- | -------- | --------------------------------------------------------------------------------------------------------------- |
| `hostname`           | Yes      | Public hostname, e.g. `rss.mydomain.com`, no `https://` or trailing path, the add-on refuses to start otherwise |
| `oidc.issuer`        | Yes      | Your OIDC provider's issuer URL                                                                                 |
| `oidc.client_id`     | Yes      | OIDC client ID                                                                                                  |
| `oidc.client_secret` | Yes      | OIDC client secret                                                                                              |
| `oidc.provider_name` | No       | Display name on the login button                                                                                |
| `base_path`          | No       | Mount yarr under a URL sub-path (e.g. `/rss`) instead of the root                                               |
| `fever_password`     | No       | Password for Fever-compatible mobile clients, leave unset to disable Fever access (see [Mobile clients](#mobile-clients-fever-api)) |
| `debug`              | No       | Enable verbose diagnostic logging                                                                               |

Your OIDC provider must support standard auto-discovery.

## Ports

Only `3000/tcp` matters, point your reverse proxy/tunnel at it (e.g. `http://{SLUG}-yarr:3000`). It's not exposed to the host by default.

## Persistent data

| Path                  | Contents                                                                                                   |
| --------------------- | ---------------------------------------------------------------------------------------------------------- |
| `/data/yarr.db`       | SQLite database (feeds, articles)                                                                          |
| `/data/cookie_secret` | Login session signing key, generated once on first run, losing it just logs everyone out                   |

## Mobile clients (Fever API)

Set `fever_password` to enable Fever-compatible clients (Reeder, Unread, FeedMe), which authenticate per-request with their own username/password instead of OIDC; the username is always `fever`, the password is whatever you set. Enter the hostname, `fever`, and that password into your client's Fever setup screen. Leave `fever_password` unset if you don't need mobile clients — the web UI works the same either way.

## Limitations

- **No multi-user support**: everyone who logs in shares the same feed list, read/unread state, and settings. For separate feed lists per person, run separate instances.

# Fusion

Self-hosted [Fusion](https://github.com/0x2E/fusion), a lightweight RSS/Atom reader with a bundled SQLite database and native OIDC login, no extra proxy needed.

## Prerequisites

- A public URL pointing at your Home Assistant (e.g. via Cloudflare tunnel, Nginx reverse proxy, etc.) with HTTPS in front.
- An OIDC provider with a client created for this instance (e.g. [Pocket ID](https://pocket-id.org)), redirect URI as described under [Callback URL](#callback-url) below.

## Configuration

| Option                    | Required | Description                                                                                              |
| -------------------------- | -------- | ------------------------------------------------------------------------------------------------------------ |
| `external_url`             | Yes      | Public URL of this instance, including scheme, e.g. `https://fusion.mydomain.com`                        |
| `oidc.issuer`               | Yes      | Your OIDC provider's issuer URL                                                                           |
| `oidc.client_id`            | Yes      | OIDC client ID                                                                                             |
| `oidc.client_secret`        | Yes      | OIDC client secret                                                                                         |
| `oidc.allowed_user`         | No       | Restrict login to a single identity (email or subject claim)                                              |
| `fever_username`            | No       | Username for Fever-compatible mobile clients, e.g. Reeder, Unread, FeedMe (default: `fusion`)             |
| `allow_private_feeds`       | No       | Allow pulling feed URLs on private/localhost networks                                                     |
| `debug`                     | No       | Enable verbose diagnostic logging                                                                          |

## Callback URL

Your OIDC client's redirect URI must be:

```
<external_url>/api/oidc/callback
```

## Ports

Only `8080/tcp` matters, point your reverse proxy/tunnel at it (e.g. `http://{SLUG}-fusion:8080`). It's not exposed to the host by default.

## Persistent data

| Path                     | Contents                                                                                |
| ------------------------- | ---------------------------------------------------------------------------------------- |
| `/data/fusion.db`         | SQLite database (feeds, articles, groups)                                              |
| `/data/fusion_password`   | Random password generated on first run, unused day-to-day since login is via OIDC      |

## Mobile clients (Fever API)

Fusion's Fever API uses its own username/API key (`fever_username`), separate from the OIDC login. Point a Fever-compatible client (Reeder, Unread, FeedMe) at this instance using that username; see [Fusion's Fever API guide](https://github.com/0x2E/fusion/blob/main/docs/fever-api.md) for the exact client-side setup.

## Limitations

- **Feed filtering/rules, notifications, and AI features** aren't part of Fusion, by design.

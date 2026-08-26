# Glance

Self-hosted [Glance](https://github.com/glanceapp/glance), a fast, self-hosted dashboard that pulls in RSS feeds, subreddits, Hacker News, weather, server stats, and more into one customizable page.

Login is via OIDC only, set your provider's details below and there's nothing else to run or configure.

## Prerequisites

- A public hostname pointing at your Home Assistant (e.g. via Cloudflare tunnel, Nginx reverse proxy, etc.) with HTTPS in front.
- An OIDC provider with a client created for this instance (e.g. [Pocket ID](https://pocket-id.org)), redirect URI `https://<hostname>/oauth2/callback`.

## Configuration

| Option               | Required | Description                                                                                                     |
| --------------------- | -------- | ----------------------------------------------------------------------------------------------------------------- |
| `hostname`            | Yes      | Public hostname, e.g. `glance.mydomain.com`, no `https://` or trailing path, the add-on refuses to start otherwise |
| `oidc.issuer`         | Yes      | Your OIDC provider's issuer URL                                                                                 |
| `oidc.client_id`      | Yes      | OIDC client ID                                                                                                   |
| `oidc.client_secret`  | Yes      | OIDC client secret                                                                                               |
| `oidc.provider_name`  | No       | Display name on the login button                                                                                |
| `debug`               | No       | Enable verbose diagnostic logging                                                                               |

Your OIDC provider must support standard auto-discovery.

## Ports

Only `3000/tcp` matters, point your reverse proxy/tunnel at it (e.g. `http://{SLUG}-glance:3000`). It's not exposed to the host by default.

## Persistent data

| Path               | Contents                                                                    |
| -------------------- | -------------------------------------------------------------------------- |
| `/data/glance.yml`   | Your dashboard: pages, columns, and widgets — generated once on first run |

**Glance's own `auth:` config is intentionally left unset** — OIDC via oauth2-proxy is the only login gate, so there's no second password to manage. `glance.yml` is only generated on first run; after that it's yours to edit directly (via the file, Glance has no in-app editor) — the add-on won't overwrite it on restart.

## Limitations

- **No multi-user support**: whoever logs in sees the same dashboard.

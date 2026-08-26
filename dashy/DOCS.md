# Dashy

Self-hosted [Dashy](https://dashy.to), a customizable dashboard/homepage, with native OIDC login (Dashy supports OIDC directly, no extra auth proxy needed).

## Prerequisites

- A public URL pointing at your Home Assistant (e.g. via Cloudflare tunnel, Nginx reverse proxy, etc.) with HTTPS in front.
- An OIDC provider with a client created for this instance (e.g. [Pocket ID](https://pocket-id.org)). Dashy's login happens entirely in the browser, so create this as a **public client** (no client secret) — its redirect URI is just `<external_url>` (the site root, not a callback path).

## Configuration

| Option              | Required | Description                                                                                   |
| -------------------- | -------- | ------------------------------------------------------------------------------------------------- |
| `external_url`       | Yes      | Public URL of this instance, including scheme, e.g. `https://dashy.mydomain.com`               |
| `oidc.issuer`        | Yes      | Your OIDC provider's issuer URL (not the `.well-known` URL)                                    |
| `oidc.client_id`     | Yes      | OIDC client ID                                                                                   |
| `oidc.admin_group`   | No       | OIDC group name granted permission to edit the dashboard; without it, everyone who logs in can edit |

## Ports

Only `8080/tcp` matters, point your reverse proxy/tunnel at it (e.g. `http://{SLUG}-dashy:8080`). It's not exposed to the host by default.

## Persistent data

| Path             | Contents                                                                                        |
| ----------------- | --------------------------------------------------------------------------------------------------- |
| `/data/conf.yml`  | Your dashboard: sections, links, widgets, and the OIDC settings                                  |

Your dashboard layout (sections, links, widgets) is only written on first run — Dashy owns it after that, editing it through its UI. The OIDC options above are reapplied from the add-on's config every time it starts, so changing them here takes effect on restart.

## Limitations

- **Not confirmed to work behind a URL sub-path** — mount it on its own hostname, not e.g. `/dashy`.

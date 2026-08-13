# SplitPro

Self-hosted [SplitPro](https://github.com/oss-apps/split-pro), an open-source alternative to Splitwise for splitting expenses with friends, roommates, and groups.

Login is via OIDC only, set your provider's details below and there's nothing else to run or configure. A database is bundled, nothing external to set up.

## Prerequisites

- A public URL pointing at your Home Assistant (e.g. via Cloudflare tunnel, Nginx reverse proxy, etc.) with HTTPS in front.
- An OIDC provider with a client created for this instance (e.g. [Pocket ID](https://pocket-id.org)), redirect URI as described under [Callback URL](#callback-url) below.

## Configuration

| Option                              | Required | Description                                                                                                          |
| ------------------------------------ | -------- | ---------------------------------------------------------------------------------------------------------------------- |
| `external_url`                      | Yes      | Public URL of this instance, including scheme, e.g. `https://splitpro.mydomain.com`                                 |
| `oidc.well_known_url`               | Yes      | Your OIDC provider's discovery document URL                                                                          |
| `oidc.client_id`                    | Yes      | OIDC client ID                                                                                                       |
| `oidc.client_secret`                | Yes      | OIDC client secret                                                                                                   |
| `oidc.provider_name`                | No       | Display name on the login button, also used in the callback URL (default: `oidc`)                                   |
| `default_homepage`                  | No       | Landing page shown after login, e.g. `/home` or `/balances` (default: `/home`)                                      |
| `upload_max_file_size_mb`           | No       | Maximum receipt upload size in MB (default: `10`)                                                                    |
| `cache_cleanup.cron_rule`           | No       | When to clean up old cached bank/currency data, UTC cron syntax. Leave both cache_cleanup fields blank to disable it |
| `cache_cleanup.retention_interval`  | No       | How long cached data can go unused before cleanup, e.g. `2 days`                                                     |
| `debug`                             | No       | Enable verbose diagnostic logging                                                                                    |

## Callback URL

Your OIDC client's redirect URI must be:

```
<external_url>/api/auth/callback/<provider_name, lowercased, or "oidc" if not set>
```

## Ports

Only `3000/tcp` matters, point your reverse proxy/tunnel at it (e.g. `http://{SLUG}-splitpro:3000`). It's not exposed to the host by default.

## Persistent data

| Path                    | Contents                                                                                  |
| ------------------------ | -------------------------------------------------------------------------------------------- |
| `/data/postgres`        | Database (expenses, groups, balances, everything)                                         |
| `/data/uploads`         | Receipt attachments                                                                        |
| `/data/nextauth_secret` | Login session signing key, generated once on first run, losing it just logs everyone out |

## Limitations

- **Login is OIDC only**: no email/magic-link or Google sign-in, this add-on is built for a closed instance behind your own identity provider.
- **Bank sync (Plaid), push notifications, and non-default currency providers** aren't configurable through this add-on.

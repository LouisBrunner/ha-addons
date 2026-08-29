# Caddy

Generic HTTPS reverse proxy built on [Caddy](https://caddyserver.com), with
[`mholt/caddy-ratelimit`](https://github.com/mholt/caddy-ratelimit) for per-client rate
limiting built into the binary. Fronts any number of hostnames, each proxied to any
upstream reachable from this add-on (other add-ons by name, host services, LAN devices).

Each upstream is resolved at request time, not once at startup: if one backend is stopped
or unreachable, only requests to that site fail — every other site keeps working.

## Configuration

| Option                    | Required | Description                                                                                                |
| ------------------------- | -------- | ---------------------------------------------------------------------------------------------------------- |
| `base_domain`             | Yes      | Domain shared by every subdomain below, e.g. `mydomain.com`                                                |
| `email`                   | Yes      | Contact email registered with Let's Encrypt                                                                |
| `max_body_size_megabytes` | No       | Default request body size cap applied to every subdomain, `0` = unlimited (default: `16`)                  |
| `subdomains`              | No       | List of reverse-proxied sites, see below                                                                   |
| `hardening.*`             | No       | Baseline protections applied to every subdomain, see below                                                 |
| `local_certs`             | No       | Use Caddy's own local CA instead of Let's Encrypt/ZeroSSL, for testing without a publicly reachable domain |
| `debug`                   | No       | Enable verbose Caddy logging                                                                               |

Example:

```yaml
base_domain: mydomain.com
email: you@example.com
subdomains:
  - name: dashboard
    upstream: some-addon-slug:8080
  - name: registry
    upstream: some-registry-addon:5000
    max_body_size_megabytes: 0
```

Each `subdomains` entry is:

| Option                    | Required | Description                                                              |
| ------------------------- | -------- | ------------------------------------------------------------------------ |
| `name`                    | Yes      | Subdomain label, e.g. `dashboard` for `dashboard.mydomain.com`           |
| `upstream`                | Yes      | Address Caddy proxies to, e.g. `some-addon:8080` or `192.168.1.10:80`    |
| `max_body_size_megabytes` | No       | Override the top-level Max Body Size for this site only, `0` = unlimited |
| `no_store`                | No       | Send `Cache-Control: no-store` on every response from this site          |

Each entry is a plain reverse proxy: one hostname to one upstream, with the baseline
hardening below always applied.

Certificates are issued via HTTP-01 (port 80) or TLS-ALPN (port 443), whichever is
reachable from the internet. Set `local_certs` to skip this and use Caddy's own local CA
instead — for testing against a domain that isn't (or can't yet be) publicly reachable;
browsers and clients won't trust the certificate unless they trust that local CA too.

### Hardening

`hardening` sets protective defaults not otherwise applied automatically. Every knob below
is disabled by setting it to `0` (or an empty list) — there's no separate `enabled` flag.

- Client-facing read/write timeouts fixed at 60s each, not configurable.
- **`hardening.rate_limit.events`/`.window`** (default `500` requests per `10s`) — per-site,
  per-client rate limit, keyed on the resolved real client IP (see below). Set `events` to
  `0` to disable.
- **`hardening.hsts.max_age_seconds`/`.include_subdomains`** (default 2 years) — Caddy
  doesn't send HSTS by default, so this is the only thing enabling it. Set `max_age_seconds`
  to `0` to disable.
- **`hardening.block_scanner_probes`** (default: path traversal, `.env`, `.git`,
  `config.php`, `wp-config`) — a list of glob patterns matched against the request path; a
  match gets rejected with a 403 before reaching the upstream. This is a fixed list you can
  edit, not a WAF. Empty list disables it.
- **`hardening.permissions_policy`** (default: camera, microphone, geolocation, USB,
  Bluetooth, and other sensitive browser features) — a list of feature names denied via the
  Permissions-Policy header, replacing any value the upstream sets. Empty list omits the
  header entirely.

Always on, not configurable: `X-Content-Type-Options: nosniff`,
`Referrer-Policy: same-origin`, `Cross-Origin-Opener-Policy: same-origin`,
`Cross-Origin-Embedder-Policy: credentialless`,
a `Content-Security-Policy: frame-ancestors 'none'` (added, not replacing an upstream CSP),
`X-Powered-By` removed, and JSON access logging per site. Also includes an internal
`/healthz` endpoint on port 8080, used by this add-on's own container healthcheck.

Upstream-facing (backend) requests have no timeout — a slow upload or long-running query to
your app isn't cut off.

#### Resolving the real client IP

Traffic arriving through Cloudflare Tunnel is supported: the real client IP is read from the
`Cf-Connecting-Ip` header. For any other traffic, rate limiting falls back to the raw
connection IP.

## Ports

`443/tcp` and `443/udp` (HTTP/3) are published to the host by default. `80/tcp` exists in the
port map but is unpublished (`null`) by default — remap it if you need HTTP-01 challenges or
a plain-HTTP redirect.

## Persistent data

| Path          | Contents                                                           |
| ------------- | ------------------------------------------------------------------ |
| `/data/caddy` | Caddy's own state: issued certificates, OCSP staples, ACME account |

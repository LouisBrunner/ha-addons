# AT Protocol PDS

Self-hosted [AT Protocol](https://atproto.com) Personal Data Server, based on the official [bluesky-social/pds](https://github.com/bluesky-social/pds) image.

## Prerequisites

- A public domain pointing to your Home Assistant then HTTPS routing using Cloudflare tunnel, Nginx reverse proxy, etc.

## Configuration

| Option                                 | Required                                | Description                                                                                                                                                                              |
| -------------------------------------- | --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `hostname`                             | Yes                                     | Public hostname of your PDS, e.g. `pds.mydomain.com`                                                                                                                                     |
| `admin_password`                       | Yes                                     | Password for the PDS admin API                                                                                                                                                           |
| `smtp_url`                             | Yes                                     | URL of the SMTP server to use for sending emails, e.g. `smtps://account:password@provider:port`                                                                                          |
| `smtp_from`                            | No                                      | Email address to use as the sender of emails sent by the SMTP server (default: `no-reply@hostname`)                                                                                      |
| `invite_required`                      | No                                      | Require invite codes to create accounts (default: `true`)                                                                                                                                |
| `recovery_did_key`                     | No                                      | Recommended. A public key that is added to the users DID, the corresponding private key can be used to recover user accounts in case of a catastrophic loss of data. [^recovery_did_key] |
| `gatekeep.enabled`                     | No                                      | Recommended. Run [Gatekeeper](https://tangled.org/pds.dad/pds-gatekeeper) in front of your auth endpoints, enabling MFA, captcha, etc. (default: `true`)                                 |
| `gatekeep.captcha.enabled`             | No                                      | Add captcha to Gatekeeper (default: `false`)                                                                                                                                             |
| `gatekeep.captcha.hcaptcha_site_key`   | If `gatekeep.captcha.enabled` is `true` | Site Key from [hCaptcha](https://www.hcaptcha.com/)                                                                                                                                      |
| `gatekeep.captcha.hcaptcha_secret_key` | If `gatekeep.captcha.enabled` is `true` | Secret Key from [hCaptcha](https://www.hcaptcha.com/)                                                                                                                                    |
| `gatekeep.only_migrations`             | No                                      | Disable account creations on Gatekeeper (default: `false`)                                                                                                                               |
| `debug`                                | No                                      | Enable extra diagnostic logging (default: `false`)                                                                                                                                       |

[^recovery_did_key]: Can be generated using `goat key generate`, format `did:key:` expected, see [reference](https://atproto.com/guides/going-to-production#plc-key-management#plc-key-management)

## Ports

| Port       | Direction | Description                  |
| ---------- | --------- | ---------------------------- |
| `3000/tcp` | inbound   | HTTP — point your proxy here |

HTTP is not exposed to the host by default as it is assumed that you will be using a HTTPS reverse proxy to access it (e.g. `http://{SLUG}-atproto-pds:3000`).

## Persistent data

The following paths are stored in `/data` and survive app restarts and updates:

| Path                | Contents     |
| ------------------- | ------------ |
| `/data/pds/`        | PDS storage  |
| `/data/blobs/`      | Blob storage |
| `/data/secrets.env` | Secrets      |

### Secrets (auto-generated)

On first start the app generates two secrets and saves them to `/data/secrets.env`:

- `PDS_JWT_SECRET`: signs all session tokens
- `PDS_PLC_ROTATION_KEY_K256_PRIVATE_KEY_HEX`: secp256k1 key used to rotate your DID document

## Environment variables (reference)

These are set by the app and cannot be overridden via the UI at the moment:

| Variable                  | Value                              | Notes                        |
| ------------------------- | ---------------------------------- | ---------------------------- |
| `PDS_BLOB_UPLOAD_LIMIT`   | `104857600`                        | 100 MB max blob size         |
| `PDS_DID_PLC_URL`         | `https://plc.directory`            | DID resolution               |
| `PDS_BSKY_APP_VIEW_URL`   | `https://api.bsky.app`             | App view for Bluesky clients |
| `PDS_BSKY_APP_VIEW_DID`   | `did:web:api.bsky.app`             |                              |
| `PDS_REPORT_SERVICE_URL`  | `https://mod.bsky.app`             | Moderation reports           |
| `PDS_REPORT_SERVICE_DID`  | `did:plc:ar7c4by46qjdydhdevvrndac` |                              |
| `PDS_CRAWLERS`            | `https://bsky.network`             | Firehose crawlers to notify  |
| `PDS_RATE_LIMITS_ENABLED` | `true`                             |                              |

There are also a few unset ones which can be added in the future:

| Variable                               | Notes |
| -------------------------------------- | ----- |
| `GATEKEEPER_CREATE_ACCOUNT_PER_SECOND` |       |
| `GATEKEEPER_CREATE_ACCOUNT_BURST`      |       |
| `PDS_PRIVACY_POLICY_URL`               |       |
| `PDS_TERMS_OF_SERVICE_URL`             |       |
| `PDS_CONTACT_EMAIL_ADDRESS`            |       |

## First-time account creation

```sh
goat pds admin --admin-password YOUR_ADMIN_PASSWORD --pds-host https://pds.mydomain.com account create --email you@example.com --handle you.pds.mydomain.com --password yourpassword
# or
pdsadmin account create you@example.com you.pds.mydomain.com
# or (INVITE_CODE is required if `invite_required` is `true`, otherwise omit it)
curl -X POST https://pds.mydomain.com/xrpc/com.atproto.server.createAccount \
  -H "Content-Type: application/json" \
  -d '{
    "email": "you@example.com",
    "handle": "you.pds.mydomain.com",
    "password": "yourpassword",
    "inviteCode": "INVITE_CODE"
  }'
```

If using `curl` directly and `invite_required` is `true`, first generate an invite code via the admin API:

```sh
curl -X POST https://pds.mydomain.com/xrpc/com.atproto.server.createInviteCode \
  -u "admin:YOUR_ADMIN_PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{"useCount": 1}'
```

> [!NOTE]
> If you want to use [`pdsadmin`](https://github.com/bluesky-social/pds/blob/main/pdsadmin.sh),
> you will need to create a `/pds/pds.env` file (or override `PDS_ENV_FILE`) and set `export PDS_HOSTNAME=pds.mydomain.com` and `export PDS_ADMIN_PASSWORD=yourpassword` inside.

## Custom domain as handle

To use `mydomain.com` as your handle instead of `you.pds.mydomain.com`, add a DNS TXT record:

```
_atproto.mydomain.com  TXT  "did=did:plc:yourDIDhere"
```

Or serve your DID at `https://mydomain.com/.well-known/atproto-did` (plain text response).

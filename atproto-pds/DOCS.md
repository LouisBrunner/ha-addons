# AT Protocol PDS

Self-hosted [AT Protocol](https://atproto.com) Personal Data Server, based on the official [bluesky-social/pds](https://github.com/bluesky-social/pds) image.

## Prerequisites

- A public domain pointing to your Home Assistant then HTTPS routing using Cloudflare tunnel, Nginx reverse proxy, etc.

## Configuration

| Option            | Required | Description                                               |
| ----------------- | -------- | --------------------------------------------------------- |
| `hostname`        | Yes      | Public hostname of your PDS, e.g. `pds.mydomain.com`      |
| `admin_password`  | Yes      | Password for the PDS admin API                            |
| `invite_required` | No       | Require invite codes to create accounts (default: `true`) |

## Ports

| Port       | Direction | Description                  |
| ---------- | --------- | ---------------------------- |
| `3000/tcp` | inbound   | HTTP — point your proxy here |

HTTP is not exposed to the host by default as it is assumed that you will be using a HTTPS reverse proxy to access it (e.g. `http://{SLUG}-atproto-pds:3000`).

## Persistent data

The following paths are stored in `/data` and survive addon restarts and updates:

| Path                | Contents     |
| ------------------- | ------------ |
| `/data/`            | Data storage |
| `/data/blobs/`      | Blob storage |
| `/data/secrets.env` | Secrets      |

### Secrets (auto-generated)

On first start the addon generates two secrets and saves them to `/data/secrets.env`:

- `PDS_JWT_SECRET`: signs all session tokens
- `PDS_PLC_ROTATION_KEY_K256_PRIVATE_KEY_HEX`: secp256k1 key used to rotate your DID document

## Environment variables (reference)

These are set by the addon and cannot be overridden via the UI at the moment:

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
> you will need to create a `/pds/pds.env` file (or override `PDS_ENV_FILE`) and set your `PDS_HOSTNAME=pds.mydomain.com` and `PDS_ADMIN_PASSWORD=yourpassword`.

## Custom domain as handle

To use `mydomain.com` as your handle instead of `you.pds.mydomain.com`, add a DNS TXT record:

```
_atproto.mydomain.com  TXT  "did=did:plc:yourDIDhere"
```

Or serve your DID at `https://mydomain.com/.well-known/atproto-did` (plain text response).

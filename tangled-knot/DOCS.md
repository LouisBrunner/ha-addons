# Tangled Knot

Self-hosted [Tangled](https://tangled.org) knot server for Git repository hosting on the AT Protocol, based on the [knot-docker](https://tangled.org/tangled.org/knot-docker) image.

## Prerequisites

- A public domain pointing to your Home Assistant then HTTPS routing using Cloudflare tunnel, Nginx reverse proxy, etc.
- A running AT Protocol PDS (the [AT Protocol PDS app](../atproto-pds) or any other PDS)
- Your DID, find it at `https://your-pds-hostname/xrpc/com.atproto.identity.resolveHandle?handle=yourhandle`

## Configuration

| Option      | Required | Description                                          |
| ----------- | -------- | ---------------------------------------------------- |
| `hostname`  | Yes      | Public hostname of the knot, e.g. `git.mydomain.com` |
| `owner_did` | Yes      | Your DID, e.g. `did:plc:xxxxxxxxxxxx`                |

## Ports

| Port                    | Purpose                  |
| ----------------------- | ------------------------ |
| `5555/tcp`              | HTTP (point your proxy)  |
| `22/tcp` (host: `2525`) | SSH (for `git` over SSH) |

HTTP is not exposed to the host by default as it is assumed that you will be using a HTTPS reverse proxy to access it (e.g. `http://{SLUG}-tangled-knot:5555`).

## Persistent data

The following paths are stored in `/data` and survive app restarts and updates:

| Path                  | Contents                                       |
| --------------------- | ---------------------------------------------- |
| `/data/knotserver.db` | SQLite database                                |
| `/data/repositories/` | Bare Git repositories                          |
| `/data/keys/`         | SSH host keys (symlinked from `/etc/ssh/keys`) |

## Registering your knot with Tangled

After starting the app, register it [here on Tangled](https://tangled.org/settings/knots).

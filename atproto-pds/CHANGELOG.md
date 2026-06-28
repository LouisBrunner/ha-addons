# Changelog

## 0.0.1

- Based on `ghcr.io/bluesky-social/pds:0.4.5009` and `tangled.org/pds.dad/pds-gatekeeper`
- Auto-generates `PDS_JWT_SECRET` and `PDS_PLC_ROTATION_KEY_K256_PRIVATE_KEY_HEX` on first run and persists them to `/data/secrets.env`
- Supports most options of both images and custom Caddy configuration

# RTSP Stream Fixer

Proxy that fixes non-standard RTSP streams so they work with FFmpeg and the Generic Camera integration.

## Configuration

| Option    | Required | Description                             |
| --------- | -------- | --------------------------------------- |
| `streams` | Yes      | List of streams to proxy (see below)    |
| `debug`   | No       | Enable debug logging (default: `false`) |

Each entry in `streams` accepts:

| Field                        | Required | Description                                                                                                                                                        |
| ---------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `name`                       | Yes      | Identifier used in the proxy URL, e.g. `living-room`                                                                                                               |
| `url`                        | Yes      | URL of the upstream RTSP stream, e.g. `rtsp://192.168.1.10:554/stream`                                                                                             |
| `fix_force_tcp_in_transport` | No       | Force TCP in the Transport header — fixes streams that cause `Nonmatching transport in server reply` (HA reports it as `Invalid data found when processing input`) |

## Ports

| Port       | Direction | Description       |
| ---------- | --------- | ----------------- |
| `6666/tcp` | inbound   | RTSP proxy (TCP)  |
| `6666/udp` | inbound   | RTSP proxy (UDP)  |
| `6667/tcp` | inbound   | Thumbnails (HTTP) |

By default, no ports are exposed on the host network, you will need to use `{SLUG}-rtsp-stream-fixer:6666/6667` inside your Home Assistant instance to access them.

## Usage

Point your Generic Camera integration or FFmpeg at the proxy instead of the upstream stream directly:

```
rtsp://{SLUG}-rtsp-stream-fixer:6666/NAME
```

Where `NAME` matches the `name` field of the stream in your configuration.

Thumbnails captured by the proxy are available at:

```
http://{SLUG}-rtsp-stream-fixer:6667/NAME
```

### Example

```yaml
streams:
  - name: living-room
    url: rtsp://192.168.1.10:554/stream
    fix_force_tcp_in_transport: true
```

Proxy URL: `rtsp://{SLUG}-rtsp-stream-fixer:6666/living-room`

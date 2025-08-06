# RTSP Stream Fixer

Instead of directing FFMpeg (or the Generic Camera integration) to your stream directly, you instead point it to this addon which then takes care of fixing the stream for you.

For example, if you have a stream available at rtsp://192.168.1.10:554/stream which needs fixing, you can add the following configuration to this addon:

```yaml
- name: first # this can be anything you want
  url: rtsp://192.168.1.10:554/stream
  # fixes issues where the Transport header is wrong and FFMpeg complains with "Nonmatching transport in server reply"
  # (HA usually reports it with "Invalid data found when processing input")
  fix_force_tcp_in_transport: true
```

You can then use `rtsp://YOUR_HOMEASSISTANT_HOST:6666/first` (note that `first` needs to match the name of your stream) in your integration and the server will automatically fix your stream so it works for HomeAssistant!

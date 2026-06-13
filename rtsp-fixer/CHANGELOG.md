## 0.0.5

Fix CI image naming: fork upstream prepare/publish actions to restore the original {image}-{arch} per-arch image naming convention.

## 0.0.4

Fix SIGSEGV crash when saving H264 thumbnails caused by stride-padded YCbCr buffers from openh264.

## 0.0.3

Fix servers binding to 127.0.0.1, which prevented HA addon port mapping from working.

## 0.0.2

Serves the latest recorded frame when the stream is not available. This reduces the logging in Home Assistant when your camera is often off.

## 0.0.1

Initial version

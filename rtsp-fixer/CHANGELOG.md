# Releases

## 0.0.11

Mask `streams[].url` in the UI/logs by typing it as `password` (RTSP URLs commonly embed camera credentials).

## 0.0.10

Fix `run.sh` missing from the image, causing the addon to fail to start.

## 0.0.9

Do not expose any ports by default, use `{SLUG}-rtsp-stream-fixer:6666/6667` to access the stream/thumbnails.

## 0.0.8

Add POST endpoint on the thumbnail URI to reset the capture timeout, forcing a new screenshot on the next frame.

## 0.0.7

Fix race condition when saving thumbnails and properly give timestamps in the thumbnail stream (ffmpeg compatibility).

## 0.0.6

Improve offline detection for streams and serve thumbnails with error feedback.

## 0.0.5

Fix addon update not being detected due to a change of how upstream actions build the Docker image name.

## 0.0.4

Fix SIGSEGV crash when saving H264 thumbnails caused by stride-padded YCbCr buffers from openh264.

## 0.0.3

Fix servers binding to 127.0.0.1, which prevented HA addon port mapping from working.

## 0.0.2

Serves the latest recorded frame when the stream is not available.
This reduces the logging in Home Assistant when your camera is often off.

## 0.0.1

Initial version

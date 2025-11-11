#!/bin/sh
set -e

if ! getent group sacloud-otelcol >/dev/null; then
  groupadd --system sacloud-otelcol
fi

if ! getent passwd sacloud-otelcol >/dev/null; then
  useradd --system --gid sacloud-otelcol --shell /sbin/nologin --home-dir /var/lib/sacloud-otel-collector --no-create-home sacloud-otelcol
fi

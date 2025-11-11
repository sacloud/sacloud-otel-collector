#!/bin/sh
set -e

systemctl stop sacloud-otel-collector || true
systemctl disable sacloud-otel-collector || true

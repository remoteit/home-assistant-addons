#!/bin/sh
# The add-on's start: its options (/data/options.json) into the device image's environment, the panel beside it, and
# then the image's own entrypoint (remoteit-device container/entrypoint), which registers, keeps the key in /data and
# runs connectd.
set -u

# The key lives in the add-on's /data, which Home Assistant keeps and backs up. A code entered in the panel waits in
# /data/registration_code until the device registers; one refused is not tried again, and the panel shows why.
export R3_DEVICE_STATE=/data/remoteit-device
export R3_REGISTRATION_CODE_FILE=/data/registration_code
export R3_REGISTRATION_RETRY=1
# The options, the device's policy, and Home Assistant's own UI (Supervisor's /core/info): remoteit-panel env.
eval "$(/usr/lib/remoteit-device/remoteit-panel env /data/options.json)"
: "${HA_PORT:=8123}" "${HA_SCHEME:=http}" "${R3_DEVICE_STAGE:=@@STAGE@@}"
echo "remoteit-device: Home Assistant's UI is $HA_SCHEME on 127.0.0.1:$HA_PORT"

# The panel, through ingress: Supervisor's proxy (172.30.32.2) reaches a host-network add-on at the hassio network's
# gateway, 172.30.32.1, so the panel listens there alone when the host has it — not on the LAN — and answers only the
# proxy wherever it listens. Started again should it end.
listen=
ip -4 -o addr show 2>/dev/null | grep -q ' 172\.30\.32\.1/' && listen=172.30.32.1
(
  while :; do
    /usr/lib/remoteit-device/remoteit-panel -listen "$listen:${R3_PANEL_PORT:-29190}" -state "$R3_DEVICE_STATE" -code-file "$R3_REGISTRATION_CODE_FILE" \
      -stage "$R3_DEVICE_STAGE" -ha-port "$HA_PORT" -ha-scheme "$HA_SCHEME" ${R3_PANEL_ALLOW:+-allow "$R3_PANEL_ALLOW"}
    sleep 5
  done
) &

exec /usr/lib/remoteit-device/entrypoint

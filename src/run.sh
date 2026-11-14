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

# console_admin: the console's own accounts (group r3-admins, made by the console app as people and AIs sign in) may
# manage Home Assistant — sudo, and `ha` with Supervisor's token, kept here for root alone. Off, they get neither.
# The add-on's container starts afresh each time, so nothing is left from a start with it on.
admin=/run/remoteit-addon
rm -rf "$admin" /etc/sudoers.d/remoteit-console
if [ "$(jq -r '.console_admin // false' /data/options.json 2>/dev/null)" = true ]; then
  mkdir -m 0700 "$admin"
  if [ -n "${SUPERVISOR_TOKEN:-}" ]; then
    (umask 077 && printf '%s' "$SUPERVISOR_TOKEN" > "$admin/supervisor_token")
  else
    echo "remoteit-device: console_admin: Supervisor gave no token, so ha cannot reach it"
  fi
  printf '%%r3-admins ALL=(ALL) NOPASSWD: ALL\n' > /etc/sudoers.d/remoteit-console
  chmod 0440 /etc/sudoers.d/remoteit-console
  echo "remoteit-device: console_admin: the console's accounts may use sudo and ha"
fi

# auto_update (on unless turned off): the device takes remote.it's releases itself, as a package install does — each
# into /data, handed over to beside the one running, confirmed, or the one before kept (the image's
# R3_DEVICE_SELF_UPDATE). An add-on update newer than what the device took runs as soon as it starts. Off, the device
# is updated with the add-on alone.
if [ "$(jq -r 'if .auto_update == false then "off" else "on" end' /data/options.json 2>/dev/null || echo on)" = on ]; then
  export R3_DEVICE_SELF_UPDATE=1
  echo "remoteit-device: auto_update: the device takes remote.it's releases itself"
fi
exec /usr/lib/remoteit-device/entrypoint

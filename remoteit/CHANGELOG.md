# Changelog

## 0.1.0

- First version, for the solo stage: the remote.it device (`ghcr.io/remoteit/device:1.1.0.20261009041129-solo`,
  connectd-go 5.6.1.20261009041129) on the host's network, its key in the add-on's `/data`.
- Registers with a code from the options or entered in the panel; a code refused is shown in the panel, not tried
  again.
- The panel (ingress): the device's name, state, UID, owner, versions and claim code; a link to it in the portal; the
  service to add for Home Assistant's UI (HTTP or HTTPS, as Supervisor says Home Assistant serves it).
- Off: the console, the proxy role, the subnet and exit node (no tun device), remote updates (the add-on is updated
  instead); any port off unless the option is on.

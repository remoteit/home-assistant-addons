# Changelog

## 0.3.0

- **A new add-on identity: `remoteit_solo`, "remote.it (solo test build)"** (was `remoteit`, "remote.it (solo)"), in the
  repository's `remoteit-solo/` folder — each remote.it stage has its own add-on now, and `remoteit` is kept for the
  supported remote.it release. Home Assistant's stage `experimental`: the store lists it only with Advanced mode on.
- Home Assistant keeps an add-on's data by its slug, so an install of `remoteit` is not updated to this one: install
  **remote.it (solo test build)**, and remove `remoteit`. To keep the same remote.it device (its UID and key) rather
  than register a new one, copy the old add-on's `/data/remoteit-device/` into the new one's before its first start —
  from the host (the SSH add-on with protection mode off, or the console):
  `cp -a /mnt/data/supervisor/addons/data/<old slug>/remoteit-device /mnt/data/supervisor/addons/data/<new slug>/`,
  where a slug is the add-on's as Supervisor names it (`local_remoteit` for a local add-on; `<repository hash>_remoteit_solo`
  from the repository — the add-on's Info page URL shows it). Do this with the new add-on stopped, and only one of the
  two may run with that key.
- The repository is public on GitHub: its URL needs no token.

## 0.2.0

- Built on the Home Assistant machine by Supervisor (no `image`): Home Assistant's base image
  (`ghcr.io/home-assistant/{aarch64,amd64}-base:3.24-2026.10.0`), the remote.it device 1.1.0.20261009062231 from solo's
  downloads — its manifest checked against solo's release key and the package against the manifest's SHA-256 — and the
  panel built from source. No private image and no registry login.

## 0.1.0

- First version, for the solo stage: the remote.it device (`ghcr.io/remoteit/device:1.1.0.20261009041129-solo`,
  connectd-go 5.6.1.20261009041129) on the host's network, its key in the add-on's `/data`.
- Registers with a code from the options or entered in the panel; a code refused is shown in the panel, not tried
  again.
- The panel (ingress): the device's name, state, UID, owner, versions and claim code; a link to it in the portal; the
  service to add for Home Assistant's UI (HTTP or HTTPS, as Supervisor says Home Assistant serves it).
- Off: the console, the proxy role, the subnet and exit node (no tun device), remote updates (the add-on is updated
  instead); any port off unless the option is on.

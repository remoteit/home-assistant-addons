# Changelog

## 0.13.0

- `auto_update`, on by default: the device takes remote.it's releases itself, as soon as they are out, instead of
  waiting for an add-on update. Each new version starts beside the running one and takes over only once it works.
  Turn it off to update the device with the add-on alone. An add-on update newer than what the device took always
  wins.
- The add-on runs without Docker's init (`init: false`): the device image's entrypoint is PID 1, as in the image itself,
  which its self-update needs.

## 0.12.3

- The remote.it device 1.1.0.20261114044939 (connectd-go 5.6.1.20261114044939, connectd-source e9b25cf3f4de): solo's
  current release, which the add-on follows on its own.

## 0.12.2

- The remote.it device 1.1.0.20261113184118 (connectd-go 5.6.1.20261113184118, connectd-source 0812045298fa): solo's
  current release, which the add-on follows on its own.

## 0.12.1

- The remote.it device 1.1.0.20261113022007 (connectd-go 5.6.1.20261113022007, connectd-source 74b1b158c6b9): solo's
  current release, which the add-on follows on its own.

## 0.12.0

- The device can use an exit: choose one for it in the portal, and the Home Assistant machine's internet traffic — Home
  Assistant's, Supervisor's and every add-on's, as the add-on runs on the host's network — goes out through it, your
  home network's own addresses staying direct. While the exit is down the traffic is held, not sent the usual way.
- The add-on now asks Home Assistant for `NET_ADMIN` and the tun device (`/dev/net/tun`) — the capability and the
  device alone, not full hardware access — and carries iproute2's `ip`, which an exit's routes need (the base image's
  busybox `ip` refuses them). `NET_ADMIN` lowers the add-on's security rating from 6 to 5. With them the portal may
  also turn on the device's subnet or make the machine an exit for others; both stay off until you do.
- The remote.it device 1.1.0.20261112045329 (connectd-go 5.6.1.20261112045329): an exit whose routes cannot be set is
  said to have failed, and why, rather than shown as up while the traffic goes out directly.

## 0.11.0

- The console can manage Home Assistant: the new `console_admin` option (off) gives the console's own accounts — an
  AI's commands — passwordless sudo, Home Assistant's configuration at `/config`, and `ha`, Home Assistant's CLI
  5.6.0 (`ha core restart`, `ha core logs`, `ha core check`).
- The add-on is given Home Assistant's configuration, writable, at `/config`, and Supervisor's `homeassistant` role
  (Home Assistant's core, not other add-ons or the host). The portal's root terminal reaches both whether or not
  `console_admin` is on.
- The docs say who the console runs as: the portal's terminal as root, an AI's commands as your own account.

## 0.10.3

- The remote.it device 1.1.0.20261108005210: commands an AI runs through the console as your own account (not root)
  work — they were refused "permission denied" on the console's socket, as the device's data directory let no one but
  root through.

## 0.10.2

- The console can be turned on: the device now finds the release keys 0.10.1 laid in, which it checks the console's
  download against. Under 0.10.0 and 0.10.1 turning it on was refused ("this install carries no release key").

## 0.10.1

- The console can be turned on: the device now has remote.it's release keys, which it checks the console's download
  against. Under 0.10.0 turning it on was refused ("this install carries no release key").

## 0.10.0

- The console is on: whoever manages the device can open a terminal on it from the portal, and an AI given the device
  can run commands there — a shell in the add-on's container, as root, on the host's network, by OpenSSH the device
  installs itself. The new `console` option (on) turns it off; while it is off the portal cannot turn it on.
- The remote.it device 1.1.0.20261107235831 (connectd-go 5.6.1.20261107235831). A local port's redirects are no longer
  remembered by the browser, so a service changed from http to https (or back) is followed at once.

## 0.9.0

- The remote.it device 1.1.0.20261103215803 (connectd-go 5.6.1.20261103215803). Upgrades hand over without a gap:
  the new device takes over from the old one in seconds and the old one is told to stop, every time. A connection
  closed from the portal or an app ends at once, also when one end had already forgotten it, and a connection that
  ends on its own is reported until the service hears it. A local port can be served over https
  (`<name>.at.remote.it`), a service's redirects keep the local port, and a service set up as http that answers only
  with a redirect to https is followed to https and reported, so the portal offers to switch it. The device speaks the
  remote.it service's current protocol, so what it finds about itself reaches the portal again.

## 0.8.0

- The remote.it device 1.1.0.20261028223058 (connectd-go 5.6.1.20261028223058): idle connections cost far less CPU, and
  setting up a connection takes less work, so a machine holding many connections stays light. A connection whose keys
  are renewed keeps going rather than starting again, and an IPv6 path that was slow to answer is tried once more. With
  the console turned on, a long command an AI starts keeps running on the device by itself, its output kept there to be
  read, rather than ending with the connection it was started over.

## 0.7.0

- The remote.it device 1.1.0.20261021193149 (connectd-go 5.6.1.20261021193149): each connection to a service now
  carries where it goes, so a service you add or change in the portal works at once, with nothing to wait for on the
  device. In the portal, **What this device sees** lists what your Home Assistant machine finds on your network (printers,
  file shares, web pages and other devices that announce themselves) and what runs on the machine itself, each with
  **Add as service**. A service at a `.local` name follows the device as its address changes, and a connection to a port
  nothing listens on is refused at once with the reason in the log. Printers are ordinary services now (type Printer):
  add one from **What this device sees**. An AirPrint service added under 0.5.0 is a Printer service now: open it in
  the portal and save it once, and the Macs that reach it list it in their print dialogs again.

## 0.6.0

- The remote.it device 1.1.0.20261010045250 (connectd-go 5.6.1.20261010045250): a session carried by the other
  device's reflector keeps looking for a direct path every 30 s, as the relaying side does, so one that becomes
  possible is found; a stalled connection to the reflector no longer freezes the device (relayed sends are queued and
  a write stuck past 5 s redials).

## 0.5.0

- The remote.it device 1.1.0.20261010000855 (connectd-go 5.6.1.20261010000855): an **AirPrint** service on this device
  (type AirPrint, the printer's address, port 631) makes its printer appear in the print dialogs of the Macs that reach
  it — the device reads the printer's own Bonjour record and passes it on. And the scheme a web service really answers
  (HTTP or HTTPS) now reaches the portal, which offers to change the service's type to match.

## 0.4.0

- The panel, and the device image's entrypoint and healthcheck, come with the remote.it device from solo's downloads
  (the build's `container-<arch>` package, checked against the same signed manifest) instead of being built or copied
  here: the install no longer builds Go, so it is quicker. Nothing changes in what the add-on does.
- The remote.it device 1.1.0.20261009223419 (connectd-go 5.6.1.20261009223419).

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

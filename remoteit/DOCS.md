# remote.it

This add-on makes the machine Home Assistant runs on a remote.it device. Through it you reach Home Assistant's own
web UI from anywhere, and — because the add-on runs on the host's network — any other machine on your home network:
your router, a camera, a NAS. Nothing is opened on your router: the device connects out to remote.it.

It is for remote.it's **solo** stage for now: register it with a code from <https://app.solo.remote.it>.

## Installing from the private repository

The repository and its image are private while the add-on is tested, so Home Assistant needs a GitHub login for each.

1. **A GitHub token** with read access to both:
   - for the repository: a fine-grained token on `remoteit/home-assistant-addons` with *Contents: read*, or a classic
     token with `repo`;
   - for the image in GitHub's registry: a **classic** token with `read:packages` (GitHub's registry does not take
     fine-grained tokens).

   One classic token with `repo` and `read:packages` serves both.
2. **The registry login** — Settings → Add-ons → Add-on Store → ⋮ → **Registries** → Add registry:
   registry `ghcr.io`, user name your GitHub user, password the token. (Or Supervisor's API:
   `POST /docker/registries` with `{"ghcr.io": {"username": "<user>", "password": "<token>"}}`.) Without it the
   install fails pulling `ghcr.io/remoteit/home-assistant-remoteit`.
3. **The repository** — Settings → Add-ons → Add-on Store → ⋮ → **Repositories** → add
   `https://<user>:<token>@github.com/remoteit/home-assistant-addons`. Supervisor clones it with that login (and keeps
   the URL, token and all, in its configuration).
4. Install **remote.it (solo)** from the store, then start it. It shows in the sidebar as **remote.it**.

### Or as a local add-on (testing)

Copy the `remoteit/` directory to `/addons/remoteit` on the Home Assistant machine (the Samba or SSH add-on), delete the
`image:` line from its `config.yaml`, and install it from **Local add-ons**: Supervisor then builds it on the machine
from its Dockerfile — which still pulls the private base image, so the registry login (step 2) is needed all the same.

## Registering the device

The device registers once, with a code; after that it signs in with its own key, kept in the add-on's data (and so in
Home Assistant's backups).

1. In the portal, **Add device**, and copy the registration code it shows. A code made with a service for Home Assistant
   (below) gives you its UI the moment the device registers.
2. Paste it in the **remote.it panel** and select Register — or put it in the add-on's `registration_code` option and
   restart the add-on.

The panel then shows the device's name, state and UID, and **Open in remote.it** takes you to it in the portal. A code
that is refused is shown in the panel with the reason, and not tried again: enter another. If the device lands
*unclaimed* (a code that registers devices without an owner), the panel shows its claim code, to claim it in the
portal.

## Home Assistant's web UI

Home Assistant listens on this machine, and the add-on runs on the host's network, so the service for its UI is on
**127.0.0.1**, at Home Assistant's port (8123 unless you changed it). In the portal, on this device, **Add service**:

| | |
|---|---|
| Type | **HTTP** — or **HTTPS** if Home Assistant serves TLS itself (`ssl_certificate` in its `http:` configuration) |
| Host | `127.0.0.1` |
| Port | `8123` |

The panel shows the type and port to use: it asks Supervisor how Home Assistant serves its UI. The type matters: an
HTTP service in front of a port that answers TLS gets an empty reply.

Home Assistant sees these connections come from 127.0.0.1, with no forwarding headers, so it needs no
`trusted_proxies` setting for them.

## Your home network: a gateway

Any service on this device may point at another machine on your network: in the portal, **Add service** with that
machine's address as its host — your router at `192.168.1.1` port 80 or 443, a camera's RTSP port, a NAS's web UI. The
device connects to it from the Home Assistant machine, so whatever that machine reaches, you reach.

The device does nothing more than the services you add: **any port** (connections to any host and port it reaches,
without a service for each) is off unless you turn on the `any_port` option.

## Options

| Option | |
|---|---|
| `registration_code` | a registration code, used once on the first start (or enter it in the panel) |
| `device_name` | the name the device registers with, unless the code names it (`Home Assistant`) |
| `stage` | remote.it's environment: `solo` |
| `any_port` | let the portal open connections to any port on any host this machine reaches (off) |

## What the device does not do here

- **No console**: remote.it's terminal would be a shell in this add-on's container, on the host's network.
- **No subnet and no exit node**: they need a network device the add-on is not given.
- **No proxy role.**
- **No remote updates**: the device is updated with the add-on (the portal shows its updates as off).

These are the device's local policy, which the portal cannot change.

## Backups and a second machine

The device's key is in the add-on's data and in every Home Assistant backup. Restoring a backup on a second machine
while the first still runs gives both the same device: they take it from each other at each sign-in. Restore onto a
machine that replaces the first, or remove the add-on's data on one of them and register it again.

## The panel and your network

The panel is served only through Home Assistant's sidebar (ingress): it answers Supervisor's proxy (172.30.32.2) alone,
and listens on the hassio network's address on the host (172.30.32.1:29190), not on your LAN. It holds no remote.it
credential and calls no remote.it API: what it shows is the device's own account.

# remote.it Home Assistant add-ons

A Home Assistant add-on repository with one add-on, **remote.it** ([`remoteit/`](remoteit/)): the Home Assistant
machine as a remote.it device — Home Assistant's own web UI reachable through remote.it, and the machine as a gateway
to the other devices on its home network.

Private, for remote.it's **solo** stage first. Supervisor builds the add-on on the Home Assistant machine from its
Dockerfile — Home Assistant's public base image, the remote.it device from solo's downloads
(https://downloads.solo.remote.it/device/, checked against solo's release key) and the panel from source — so it needs
no private image and no registry login.

| | |
|---|---|
| Installing it | [remoteit/DOCS.md](remoteit/DOCS.md#installing) |
| Using it | [remoteit/DOCS.md](remoteit/DOCS.md) |
| A newer device build | `DEVICE_VERSION` in `remoteit/Dockerfile` (a version in https://downloads.solo.remote.it/device/), and the add-on's version |
| Changes | [remoteit/CHANGELOG.md](remoteit/CHANGELOG.md) |

## Layout

- `repository.yaml` — what Home Assistant reads when the repository is added.
- `remoteit/config.yaml` — the add-on: `aarch64` and `amd64` (Home Assistant no longer supports armhf, armv7 and
  i386), the host's network, ingress, its options. No `image`: Supervisor builds it on the machine.
- `remoteit/build.yaml` — Home Assistant's public base image per arch, pinned.
- `remoteit/Dockerfile` — the panel (`panel/`, Go, standard library only, built with no network), the device fetched
  and checked by `device/fetch`, the device image's entrypoint and healthcheck, and `run.sh`.
- `remoteit/device/fetch` — downloads the device's Linux tarball for the arch from solo's downloads; the version's
  `manifest.json` must be signed (Ed25519) by a key in `device/release-keys`, and the tarball must match the manifest's
  SHA-256, as the device's installer (remoteit-device `install.sh`) checks a first install.
- `remoteit/device/release-keys` — solo's release public key: remoteit-device's `release-keys-solo`.
- `remoteit/device/entrypoint`, `remoteit/device/healthcheck` — copies, byte for byte, of remoteit-device's
  `container/entrypoint` and `container/healthcheck` (branch `home-assistant`, 4b240a9): that is their source. Check with
  `git -C <remoteit-device> show <rev>:container/entrypoint | diff - remoteit/device/entrypoint`.

## Checks

- Frenck's add-on linter (`frenck/action-addon-linter` v2.21.1, its `src/` run in Docker with `INPUT_PATH=remoteit`):
  passes.
- `device/fetch` refuses a manifest signed by another key, a package whose SHA-256 differs from the manifest's, and an
  arch with no package.
- Built by Supervisor on a HAOS 18.3 aarch64 VM as a local add-on: 77 s to install (pulling the base and Go images,
  fetching the device, building the panel), the image 106 MB; it signed in as the device it was before (its
  `remoteit-device/` carried over in `/data`), its panel showed it connected, and Home Assistant's UI answered 200 through
  solo's proxy.

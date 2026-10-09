# remote.it Home Assistant add-ons

A Home Assistant add-on repository with one add-on, **remote.it** ([`remoteit/`](remoteit/)): the Home Assistant
machine as a remote.it device — Home Assistant's own web UI reachable through remote.it, and the machine as a gateway
to the other devices on its home network.

Private, for remote.it's **solo** stage first. The add-on's image is `ghcr.io/remoteit/home-assistant-remoteit` (private),
built on the remote.it device image `ghcr.io/remoteit/device` (private; `container/` in
[remoteit/remoteit-device](https://github.com/remoteit/remoteit-device), docs/container-image.md).

| | |
|---|---|
| Installing it (a private repository and a private image) | [remoteit/DOCS.md](remoteit/DOCS.md#installing-from-the-private-repository) |
| Using it | [remoteit/DOCS.md](remoteit/DOCS.md) |
| Publishing an image | `BUILDX_BUILDER=<a docker-container builder> bin/publish` — every arch in config.yaml, tagged with its version |
| Changes | [remoteit/CHANGELOG.md](remoteit/CHANGELOG.md) |

## Layout

- `repository.yaml` — what Home Assistant reads when the repository is added.
- `remoteit/config.yaml` — the add-on: `aarch64` and `amd64` (Home Assistant no longer supports armhf, armv7 and
  i386), the host's network, ingress, its options. `image` names one multi-architecture image (no `{arch}`): Supervisor
  pulls it for the machine's platform.
- `remoteit/build.yaml` — the base image per arch, for Home Assistant's builder.
- `remoteit/Dockerfile` — the device image, plus the panel (`panel/`, Go, built on the build's own platform) and
  `run.sh`.

## Checks

- Frenck's add-on linter (`frenck/action-addon-linter` v2.21.1, its `src/` run in Docker with `INPUT_PATH=remoteit`):
  passes.
- Home Assistant's builder (`ghcr.io/home-assistant/aarch64-builder --aarch64 --target remoteit --test`): builds — with
  the private base served from a local registry, since the builder's container has no login to ghcr.io.

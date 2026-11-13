# remote.it Home Assistant add-ons

Home Assistant add-ons that make the Home Assistant machine a remote.it device: Home Assistant's own web UI reachable
through remote.it, and the machine as a gateway to the other devices on its home network.

## Installing

In Home Assistant: **Settings → Add-ons → Add-on Store → ⋮ → Repositories**, add

    https://github.com/remoteit/home-assistant-addons

then install the add-on from the store and start it. Its Documentation tab has the rest (registering the device,
the service for Home Assistant's UI, the gateway).

## The add-ons

One per remote.it stage:

| Add-on | Folder | remote.it stage | Home Assistant stage |
|---|---|---|---|
| remote.it (solo test build) | [`remoteit-solo/`](remoteit-solo/) | solo ([app.solo.remote.it](https://app.solo.remote.it)) | experimental |
| remote.it | `remoteit/` — coming | prod ([app.remote.it](https://app.remote.it)) | stable |

**Test builds are unsupported**: they are for testing remote.it's other stages, may change or break without notice, and
remote.it support does not cover them. Home Assistant marks them *experimental* and lists them in the Add-on Store
only with **Advanced mode** on (your user profile → Advanced mode).

**prod: coming.** The supported add-on, `remoteit/`, waits for the remote.it device package on prod's downloads
(`https://downloads.remote.it/device/`) and prod's release key (remoteit-device's `release-keys`): neither is published
yet.

## How the add-ons are built

Supervisor builds each add-on on the Home Assistant machine from its Dockerfile — no private image, no registry login:
Home Assistant's public base image, the remote.it device package downloaded from the stage's downloads and checked
against the stage's release key (`device/fetch`: the version's `manifest.json` signed with Ed25519, the tarball matching
its SHA-256, as the device's `install.sh` checks a first install), and from the same build its container pieces — the
device image's entrypoint and healthcheck and the add-on's panel, built and released with the device (remoteit-device
`container/`, manifest key `container-<arch>`). The repository holds no copy of them and builds no Go.

## One source, a folder per stage

The add-on folders are **generated**: edit `src/` and `stages/`, then run `bin/gen`. Stages may be built from the same
branch, so each has its own folder rather than its own branch.

- `src/` — the add-on, with `@@NAME@@` placeholders for a stage's values and `@@if test@@` / `@@if stable@@` …
  `@@end@@` blocks for what only a test build (or only a supported one) says.
- `stages/<stage>/stage.env` — the stage's values: the folder, name, slug, version, Home Assistant stage, remote.it
  stage, portal, downloads and device version.
- `stages/<stage>/release-keys` — the stage's release public key (remoteit-device's `release-keys-<stage>`), written to
  `device/release-keys`.
- `stages/<stage>/CHANGELOG.md` — the add-on's changelog (each stage's add-on has its own versions).
- `bin/gen` writes every stage's folder whole; `bin/gen --check` changes nothing and fails if a folder is out of date.
  CI (`.github/workflows/check.yaml`) runs it, and Frenck's add-on linter on each folder — add a new stage's folder to
  its matrix.

## Following the device's releases

The solo add-on follows stable's current on its own. Every 10 minutes `.github/workflows/follow-current.yml` runs
`bin/follow-current`, which reads `<DOWNLOADS>/channels/stable.json` (graphql writes it on every promotion; until it is
published, the top-level `<DOWNLOADS>/manifest.json` stands in). When that version is not `DEVICE_VERSION` in
`stages/solo/stage.env`, the script:

- requires the version's manifest to have the container pieces and tarballs for both of Home Assistant's arches
  (`container-arm64`, `container-amd64`, `tar-arm64`, `tar-amd64`);
- sets `DEVICE_VERSION`, bumps `ADDON_VERSION`'s patch number, and adds a CHANGELOG entry from the manifest (the
  version, connectd and its commit);
- runs `bin/gen`.

The workflow then runs `bin/gen --check`, `device/fetch` for both arches (the signature and SHA-256s Supervisor's build
checks), and the add-on linter. It then commits as github-actions and pushes to `main`, using `GITHUB_TOKEN` only.
That push starts no other workflow, so `check.yaml` does not run on it; the same checks ran before the push.
It never moves backwards: a release older than `DEVICE_VERSION` (by its date stamp) is refused.

To move by hand, either run the workflow with a `version` (**Actions → follow-current → Run workflow**), or run
`bin/follow-current [version]` locally, commit and push. To say more than the generated line, edit the new CHANGELOG
entry before pushing.
A change to the panel, the entrypoint or the healthcheck is a remoteit-device change and a new device build.

## Checks

- `bin/gen --check`, and Frenck's add-on linter (`frenck/action-app-linter` v2.21.1) on each folder — on every push
  and pull request (`check.yaml`), and in `follow-current.yml` before each bump it pushes.
- `device/fetch` refuses a manifest signed by another key, a package whose SHA-256 differs from the manifest's, and an
  arch with no package.

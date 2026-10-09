# The remote.it add-on: the remote.it device image (remoteit-device's container/, ghcr.io/remoteit/device) with the
# panel and the add-on's start script. BUILD_FROM is build.yaml's: Home Assistant's builder passes it per arch.
ARG BUILD_FROM=ghcr.io/remoteit/device:1.1.0.20261009041129-solo

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS panel
ARG TARGETARCH
ARG TARGETVARIANT
WORKDIR /src
COPY panel/ .
RUN if [ -n "$TARGETVARIANT" ]; then export GOARM=${TARGETVARIANT#v}; fi; \
    CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -trimpath -ldflags "-s -w" -o /remoteit-panel .

FROM ${BUILD_FROM}
ARG BUILD_VERSION
COPY --from=panel /remoteit-panel /usr/bin/remoteit-panel
COPY run.sh /usr/lib/remoteit-device/addon-run
# The key in the add-on's /data — set here, not only in run.sh, for the image's HEALTHCHECK to find it too.
ENV R3_DEVICE_STATE=/data/remoteit-device
LABEL io.hass.type="addon" \
      io.hass.name="remote.it" \
      io.hass.version="${BUILD_VERSION}" \
      org.opencontainers.image.title="remote.it Home Assistant add-on" \
      org.opencontainers.image.source="https://github.com/remoteit/home-assistant-addons"
ENTRYPOINT ["/usr/lib/remoteit-device/addon-run"]

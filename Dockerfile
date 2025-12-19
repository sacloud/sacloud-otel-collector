FROM gcr.io/distroless/static-debian12

ARG TARGETARCH
ARG USER_UID=10001
ARG USER_GID=10001
USER ${USER_UID}:${USER_GID}

COPY dist/sacloud-otel-collector_linux_${TARGETARCH}*/sacloud-otel-collector /sacloud-otel-collector
COPY config.yaml /etc/otelcol-contrib/config.yaml
ENTRYPOINT ["/sacloud-otel-collector"]
CMD ["--config", "/etc/otelcol-contrib/config.yaml"]
EXPOSE 4317 4318 8888 13133

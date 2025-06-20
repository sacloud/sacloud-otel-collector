FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 AS certs
RUN apk --update add ca-certificates

FROM scratch

ARG TARGETARCH
ARG USER_UID=10001
ARG USER_GID=10001
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY dist/sacloud-otel-collector_linux_${TARGETARCH}*/sacloud-otel-collector /sacloud-otel-collector
COPY config.yaml /etc/otelcol-contrib/config.yaml
ENTRYPOINT ["/sacloud-otel-collector"]
CMD ["--config", "/etc/otelcol-contrib/config.yaml"]
EXPOSE 4317 4318 8888 13133

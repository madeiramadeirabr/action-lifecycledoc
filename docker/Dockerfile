# syntax=docker/dockerfile-upstream:master-labs
FROM alpine:3.16

ADD --checksum=sha256:277baa4ac078404839e5bc565ea8c9c3013f451a870571110b0a89b3fd5e01de https://github.com/madeiramadeirabr/action-lifecycledoc/releases/download/v0.1.0/lifecycledoc_linux_amd64 /usr/bin/lifecycledoc

RUN chmod +x /usr/bin/lifecycledoc

COPY docker/entrypoint.sh entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
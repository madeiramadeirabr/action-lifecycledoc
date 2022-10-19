# syntax=docker/dockerfile-upstream:master-labs
FROM alpine:3.16

ADD https://github.com/madeiramadeirabr/action-lifecycledoc/releases/download/v0.1.0/lifecycledoc_linux_amd64 /usr/bin/lifecycledoc

RUN chmod +x /usr/bin/lifecycledoc

COPY docker/entrypoint.sh entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
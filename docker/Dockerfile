FROM alpine:3.16

ADD https://github.com/madeiramadeirabr/action-lifecycledoc/releases/download/v0.2.0/lifecycledoc_linux_amd64 /usr/bin/lifecycledoc

RUN chmod +x /usr/bin/lifecycledoc

COPY entrypoint.sh entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
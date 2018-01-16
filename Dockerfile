# Use a real distro instead of scratch in order to have some precanned certificate authorities
FROM alpine:3.7
RUN apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
ADD oidc-token-ferry /
ENTRYPOINT [ "/oidc-token-ferry" ]

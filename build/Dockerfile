# USE ALPINE LINUX AS BASE IMAGE (TO ALLOW BASH NAVIGATION)
FROM envoyproxy/envoy-alpine:3016680a511cd08d60cbcb8b8338cebb76ca7dca

MAINTAINER DKF-Basel <info@dkfbasel.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2017"

# install openssl to generate server and client certificates
# note: maybe we can switch to a go implementation in the future
RUN apk add --update openssl && \
    rm -rf /var/cache/apk/*

# create an application volume
RUN mkdir /app

# create directories for envoyproxy configuration, logs
RUN mkdir -p /app/envoy/config
RUN mkdir -p /app/envoy/log

# create a directory for certificates to be used by envoyproxy
# note: this should be mounted from the host to persist the configuration
RUN mkdir -p /app/envoy/certificates

# add the kolumbus binary
ADD kolumbus /app/bin/kolumbus

# add an envoy-configuration to use envoy as proxy for all other containers
ADD envoyproxy.config.tmpl /app/envoy/config/envoyproxy.config.tmpl

# SET THE CURRENT WORKING DIRECTORY
WORKDIR /app

# START THE APPLICATION WITH THE CONTAINER
CMD ["/app/bin/kolumbus"]

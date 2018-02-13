# USE ALPINE LINUX AS BASE IMAGE (TO ALLOW BASH NAVIGATION)
FROM envoyproxy/envoy-alpine:3016680a511cd08d60cbcb8b8338cebb76ca7dca

MAINTAINER DKF-Basel <info@dkfbasel.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2017"

# create an application volume
RUN mkdir /app

# create directories for envoy configuration and logs
RUN mkdir -p /app/envoy/config
RUN mkdir -p /app/envoy/log

# add the kolumbus binary
ADD kolumbus /app/bin/kolumbus

# add an envoy-configuration to use envoy as proxy for all other containers
ADD proxy/config.json /app/envoy/config/proxy.json

# SET THE CURRENT WORKING DIRECTORY
WORKDIR /app

# START THE APPLICATION WITH THE CONTAINER
CMD ["/app/bin/kolumbus"]

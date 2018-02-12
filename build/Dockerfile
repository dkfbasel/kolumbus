# USE ALPINE LINUX AS BASE IMAGE (TO ALLOW BASH NAVIGATION)
FROM alpine:3.7

MAINTAINER DKF-Basel <info@dkfbasel.ch>
LABEL copyright="Departement Klinische Forschung, Basel, Switzerland. 2017"

# COPY THE APPLICATIN FILES INTO THE CONTAINER
RUN mkdir /app
ADD kolumbus /app/bin/kolumbus

# SET THE CURRENT WORKING DIRECTORY
WORKDIR /app

# START THE APPLICATION WITH THE CONTAINER
CMD ["/app/bin/kolumbus"]

# Kolumbus: Creating a service mesh with docker and envoyproxy
Docker makes it easy to package your applications and run it reliably in
different environments.

However, orchestrating multiple containers with load balancing, rate limiting,
dynamic replacement of services, monitoring and all the nice operational stuff
can quickly become quite cumbersome.

Given our current move to a microservice architecture that makes heavy use of
the grpc framework, we were looking for a simple solution.

Up until now, we were managing our services manually with caddyserver as proxy
and were quite happy with this. We did look into traefik as a possible solution
but found that it would not help us with all the required aspects.

With envoyproxy and its sidecar philosophy we found a solution that seems to
work very well with our setup. One of the great things is about it is that
envoyproxy instances can fetch almost all of the configuration dynamically.

This allowed us to write a simple orchestration service ("kolumbus"), that
will watch all docker containers in the same network and use simple docker labels
(very similar to traefik) to generate the configuration required to form
an envoyproxy service mesh.

Please note that `Kolumbus` is currently pretty new and we are working on
several aspects to improve the devops experience.

An example on how it can be used is in the examples/grpc directory.

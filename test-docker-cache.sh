#!/usr/bin/env bash

cat <<DOCKERFILE > Dockerfile.deb
FROM ubuntu:24.04
RUN apt update && apt install -y \
  openssl
DOCKERFILE
docker buildx build \
  --cache-from type=gha --cache-to type=gha \
  -f Dockerfile.deb -t ubuntu:test-local .

cat <<DOCKERFILE > Dockerfile.rpm
FROM centos:7 
RUN yum install -y \
  openssl make unzip
DOCKERFILE
docker buildx build \
  --cache-from type=gha --cache-to type=gha \
  -f Dockerfile.rpm -t centos:test-local .

#docker run -t ubuntu:test-local echo hi ubuntu
#docker run -t centos:test-local echo hi centos
docker image list

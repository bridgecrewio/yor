FROM alpine:3.17.1
# checkov:skip=CKV_DOCKER_3: Not a service
# checkov:skip=CKV_DOCKER_2: Not a service
# checkov:skip=CKV2_DOCKER_2: Not a service

ARG BUILDARCH=amd64
RUN echo "linux_${BUILDARCH}"
RUN apk --no-cache add build-base git curl jq bash
RUN curl -s -k https://api.github.com/repos/bridgecrewio/yor/releases/latest  \
    | jq ".assets[] | select(.name | contains(\"linux_${BUILDARCH}\"))  \
    | .browser_download_url" -r  \
    | awk '{print "curl -L -k " $0 " -o yor.tar.gz"}' | sh
RUN tar -xf yor.tar.gz -C /usr/bin/ && rm yor.tar.gz && chmod +x /usr/bin/yor && echo 'alias yor="/usr/bin/yor"' >> ~/.bashrc
COPY entrypoint.sh /entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]

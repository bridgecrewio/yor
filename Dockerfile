FROM alpine:3.17.1
# checkov:skip=CKV_DOCKER_3: Not a service
# checkov:skip=CKV_DOCKER_2: Not a service
# checkov:skip=CKV2_DOCKER_2: Not a service

RUN apk --no-cache add build-base git curl jq bash
RUN curl -s -k https://api.github.com/repos/bridgecrewio/yor/releases/latest | jq '.assets[] | select(.name | contains("linux_386")) | select(.content_type | contains("gzip")) | .browser_download_url' -r | awk '{print "curl -L -k " $0 " -o /usr/bin/yor.tar.gz"}' | sh
RUN tar -xf /usr/bin/yor.tar.gz -C /usr/bin/ && rm /usr/bin/yor.tar.gz && chmod +x /usr/bin/yor && echo 'alias yor="/usr/bin/yor"' >> ~/.bashrc
COPY entrypoint.sh /entrypoint.sh

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -o /docker-yor-test

CMD ["/docker-yor-test"]

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/entrypoint.sh"]

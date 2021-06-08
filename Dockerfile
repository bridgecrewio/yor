FROM alpine

RUN apk update && apk --no-cache add build-base git curl jq bash
RUN curl -s https://api.github.com/repos/bridgecrewio/yor/releases/latest | jq '.assets[] | select(.name | contains("linux-386")) | select(.content_type | contains("gzip")) | .browser_download_url' -r | awk '{print "curl -L " $0 " -o /usr/bin/yor.tar.gz"}' | sh
RUN tar -xf /usr/bin/yor.tar.gz -C /usr/bin/
RUN rm /usr/bin/yor.tar.gz
RUN chmod +x /usr/bin/yor
RUN echo 'alias yor="/usr/bin/yor"' >> ~/.bashrc
COPY github_action_entrypoint.sh /github_action_entrypoint.sh

ENTRYPOINT ["/usr/bin/yor"]
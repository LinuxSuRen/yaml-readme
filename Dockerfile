FROM golang:1.17 as builder

WORKDIR /workspace
COPY . .
RUN go mod download
RUN CGO_ENABLE=0 go build -ldflags "-w -s" -o yaml-readme
RUN curl -L https://github.com/LinuxSuRen/http-downloader/releases/download/v0.0.67/hd-linux-amd64.tar.gz | tar xzv hd

FROM alpine:3.10

LABEL "com.github.actions.name"="README helper"
LABEL "com.github.actions.description"="README helper"
LABEL "com.github.actions.icon"="home"
LABEL "com.github.actions.color"="red"

LABEL "repository"="https://github.com/linuxsuren/yaml-readme"
LABEL "homepage"="https://github.com/linuxsuren/yaml-readme"
LABEL "maintainer"="Rick <linuxsuren@gmail.com>"

LABEL "Name"="README helper"

ENV LC_ALL C.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8

RUN apk add --no-cache \
        git \
        openssh-client \
        libc6-compat \
        libstdc++

COPY entrypoint.sh /entrypoint.sh
COPY --from=builder /workspace/yaml-readme /usr/bin/yaml-readme
COPY --from=builder /workspace/hd /usr/bin/hd

ENTRYPOINT ["/entrypoint.sh"]

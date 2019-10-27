FROM golang:1.13 as builder
LABEL maintainer="Sven Nebel <https://github.com/snebel29>"

ARG VERSION
ARG REPOSITORY

WORKDIR /go/src/${REPOSITORY}

# Copy everything from context to working directory within docker image
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN VERSION=${VERSION} make build

####################################################################
FROM alpine:3.9

ARG CONTAINER_USER
ARG REPOSITORY

RUN adduser -D ${CONTAINER_USER} && apk add --no-cache diffutils ca-certificates
USER ${CONTAINER_USER}
WORKDIR /
COPY --from=builder /go/src/${REPOSITORY}/kwatchman /bin

ENTRYPOINT ["kwatchman"]

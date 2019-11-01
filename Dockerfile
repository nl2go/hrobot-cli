FROM golang:alpine as build
LABEL maintainer="pahl@newsletter2go.com"
RUN apk update
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache git
WORKDIR /build
ADD . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags '-extldflags "-static"' -o hrobot-cli

FROM alpine:3.10
LABEL maintainer="pahl@newsletter2go.com"

COPY --from=build /build/hrobot-cli /hrobot-cli
ENTRYPOINT ["/hrobot-cli"]
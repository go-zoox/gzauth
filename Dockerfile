# Builder
FROM whatwewant/builder-go:v1.19-1 as builder

WORKDIR /build

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux \
  GOARCH=amd64 \
  go build \
  -trimpath \
  -ldflags '-w -s -buildid=' \
  -v -o gzauth

# Server
FROM whatwewant/go:v1.19-1

LABEL MAINTAINER="Zero<tobewhatwewant@gmail.com>"

LABEL org.opencontainers.image.source="https://github.com/go-zoox/gzauth"

ARG VERSION=latest

ENV MODE=production

COPY --from=builder /build/gzauth /bin

ENV VERSION=${VERSION}

COPY ./entrypoint.sh /

CMD /entrypoint.sh

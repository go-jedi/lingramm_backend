FROM golang:1.24.4-alpine AS builder

# install dependencies for libvips (bimg)
RUN apk add --no-cache \
    vips-dev \
    gcc \
    g++ \
    musl-dev

WORKDIR /github.com/go-jedi/lingramm_backend/app
COPY . /github.com/go-jedi/lingramm_backend/app

RUN go mod download
RUN go build -ldflags="-s -w" -trimpath -buildvcs=false -o .bin/app cmd/app/main.go

FROM alpine:latest

# install runtime-dependencies
RUN apk add --no-cache vips

WORKDIR /root/
COPY --from=builder /github.com/go-jedi/lingramm_backend/app/.bin/app .
COPY config.yaml /root/
COPY migrations /root/migrations

CMD ["./app", "--config", "config.yaml"]
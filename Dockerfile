FROM golang:1.24.0-alpine AS builder

WORKDIR /github.com/go-jedi/lingvogramm_backend/app
COPY . /github.com/go-jedi/lingvogramm_backend/app

RUN go mod download
RUN go build -ldflags="-s -w" -trimpath -buildvcs=false -o .bin/app cmd/app/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/go-jedi/lingvogramm_backend/app/.bin/app .
COPY config.yaml /root/
COPY migrations /root/migrations

CMD ["./app", "--config", "config.yaml"]
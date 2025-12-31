FROM golang:1.25.5-alpine AS builder

WORKDIR /taskpool

ENV CGO_ENABLED=0

# not sure
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-s -w" -o app ./cmd

FROM alpine:3.20

WORKDIR /taskpool

# not sure
RUN apk add --no-cache ca-certificates

COPY --from=builder /taskpool/app /taskpool/app

RUN addgroup -S app && adduser -S app -G app

USER app

EXPOSE 8080

ENTRYPOINT ["/taskpool/app"]

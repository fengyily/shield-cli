FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /shield .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /shield /usr/local/bin/shield

ENV SHIELD_LISTEN_HOST=0.0.0.0
EXPOSE 8181

ENTRYPOINT ["shield"]
CMD ["start", "--no-tray"]

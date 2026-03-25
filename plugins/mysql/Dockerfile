FROM golang:1.25 AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /shield-plugin-mysql .

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /shield-plugin-mysql /shield-plugin-mysql

EXPOSE 8080

ENV DB_HOST=127.0.0.1
ENV DB_PORT=3306
ENV DB_USER=root

ENTRYPOINT ["/shield-plugin-mysql"]

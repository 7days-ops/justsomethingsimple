FROM golang:1.25.3 AS builder

WORKDIR /app

COPY src/ .
COPY src/http/ ./http
COPY src/static ./static

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main .

FROM alpine:3.18

RUN apk add --no-cache curl

COPY --from=builder /app/main /main
COPY --from=builder /app/static /static


EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD curl -f http://localhost:8080/ || exit 1

ENTRYPOINT [ "./main" ]

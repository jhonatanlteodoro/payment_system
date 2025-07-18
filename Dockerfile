FROM golang:1.24-alpine AS builder

# Install git and timezone data
RUN apk add --no-cache git tzdata
RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY main.go ./
COPY src/ ./src/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app .

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/app /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app"]
CMD ["full-server"]
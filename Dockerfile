FROM golang:1.23-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /devices-api ./cmd/server

FROM scratch
COPY --from=builder /devices-api /devices-api
ENTRYPOINT ["/devices-api"]

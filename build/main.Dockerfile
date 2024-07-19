FROM golang:1.22.5-alpine AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./

RUN go mod download
RUN go clean --modcache

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/.bin .

EXPOSE 8080

ENTRYPOINT ["./.bin"]
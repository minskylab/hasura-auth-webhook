# first stage - builder
FROM golang:1.22 as builder

WORKDIR /app

# Copy dependency files
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy files
COPY . /app

# Build project
# CGO_ENABLED=0
RUN GOOS=linux go build -a -ldflags '-extldflags "-static"' -o app ./cmd/hasura-auth-webhook/*


# second stage - main
FROM alpine:latest

# OS dependencies
RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app .

CMD ["./app"]
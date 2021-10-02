# first stage - builder
FROM golang:1.15 as builder

WORKDIR /app

ARG GITHUB_TOKEN

ENV GO111MODULE=on
ENV GOPRIVATE="github.com/minskylab/*" 
ENV GONOSUMDB="github.com/minskylab/*"

# Copy dependency files
COPY go.mod .
COPY go.sum .

RUN git config --global url."https://x-oauth-basic:${GITHUB_TOKEN}@github.com/minskylab".insteadOf "https://github.com/minskylab"

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
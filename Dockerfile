#syntax=docker/dockerfile:1

FROM golang:1.21

# Set working directory
WORKDIR /app

# Cache go.mod for pre-downloading dependencies
COPY go.mod ./
RUN go mod download && go mod verify

# Build application
COPY . .
RUN go build -o /usr/local/bin/app ./...

EXPOSE 4000

CMD ["app"]
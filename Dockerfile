FROM golang:1.21

# Set working directory
WORKDIR /app

# Cache go.mod for pre-downloading dependencies
COPY go.mod ./
RUN go mod download && go mod verify

# Build application
COPY . .
RUN go build -o ./bin ./...

EXPOSE 4000

CMD [ "/app/bin/app" ]
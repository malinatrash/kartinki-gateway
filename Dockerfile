FROM golang:1.22.4-alpine

WORKDIR /src
# Copy proto files
COPY proto proto/

WORKDIR /src/app
# Copy gateway files
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download
COPY gateway/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go
CMD ["./main"]

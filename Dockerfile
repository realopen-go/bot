FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o realopen .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/realopen .

# Build a small image
FROM golang:alpine

COPY --from=builder /dist/realopen /

# Command to run
ENTRYPOINT ["/realopen", "run"]
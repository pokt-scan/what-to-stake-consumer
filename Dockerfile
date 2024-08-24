# Build stage
FROM golang:1.23.0 as builder
WORKDIR /app

# Download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire repository
COPY . .

# Set default build environment variables, which can be overridden at build time
ARG GOOS=linux
ARG GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}

# Build the binary and place it in /app/bin
RUN mkdir -p /app/bin && go build -ldflags '-w -extldflags "-static"' -o /app/bin/wtsc ./cmd/wtsc

# Verification step to ensure binary creation and permissions
RUN ls -la /app/bin/wtsc && chmod +x /app/bin/wtsc && ls -la /app/bin/wtsc

# Final stage
FROM alpine:latest

# Install tini and bash
RUN apk add --no-cache tini bash file

# Copy the binary from the builder stage
COPY --from=builder /app/bin/wtsc /app/bin/wtsc

# Create the results directory
RUN mkdir -p /app/results

# Ensure execute permissions explicitly
RUN chmod +x /app/bin/wtsc && ls -la /app/bin/wtsc

# Declare a volume to save results
VOLUME /app/results

WORKDIR /app

# Adjust entry point to execute the binary
CMD ["/sbin/tini", "--", "/app/bin/wtsc"]
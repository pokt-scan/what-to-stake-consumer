# Specify the base image
FROM golang:1.23.0 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
WORKDIR /app/cmd/wtsc
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o /app/wtsc .

# Start a new minimal image
FROM alpine:latest

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/wtsc /app/wtsc

# Change ownership of the files to the non-root user
RUN chown -R appuser:appgroup /app

# Declare a path for results to be mounted
VOLUME /app/results

# Use non-root user
USER appuser

# Command to run the executable
CMD ["/app/wtsc"]
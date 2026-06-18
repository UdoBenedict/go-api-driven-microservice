
# Stage 1
FROM golang:1.21-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /OneDrive/Documents/K8s/api-driven-microservice

# Copy the Go source code into the container
COPY main.go .

# Build the Go app statically. 
# CGO_ENABLED=0 ensures it doesn't depend on C libraries, making it portable.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-service main.go

# Stage 2
FROM alpine:3.22.4

# Create a dedicated non-root user and group
# -S creates a system user/group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /OneDrive/Documents/K8s/api-driven-microservice

# Copy the pre-built binary file from the builder stage
COPY --from=builder /OneDrive/Documents/K8s/api-driven-microservice/api-service .

# Change ownership of the binary to the non-root user
RUN chown -R appuser:appgroup /OneDrive/Documents/K8s/api-driven-microservice

# Tell Docker to run the container as the non-root user
USER appuser

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./api-service"]
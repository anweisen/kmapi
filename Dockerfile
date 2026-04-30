# Stage 1: Build the Go binary
FROM golang:latest as build

WORKDIR /app

# Copy the go.mod file to the container
COPY go.mod go.sum ./
COPY src/ src/
COPY static/ static/

# Generate go.sum
RUN go mod tidy

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kmapi ./src

# Build the Go application
# RUN go build -o kmapi ./src

# Stage 2: Execute the binary on a minimal Linux server image
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage to the final stage
COPY --from=build /app/kmapi .
COPY --from=build /app/static/ ./static/

EXPOSE 80

# Define the command to run when the container starts
CMD ["./kmapi"]

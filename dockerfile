# STEP 1: Use a lightweight Go base image
FROM golang:1.23.3-alpine AS builder

# STEP 2: Set working directory inside the container
WORKDIR /app

# STEP 3: Copy go.mod and go.sum files first (for caching dependencies)
COPY go.mod go.sum ./

# STEP 4: Download dependencies
RUN go mod download

# STEP 5: Copy application source code
COPY . .

# STEP 6: Build the Go application
RUN go build -o main .

# STEP 7: Use lightweight Alpine image for deployment
FROM alpine:3.18

# STEP 8: Install CA Certificates for SSL connections (required for MongoDB Atlas)
RUN apk --no-cache add ca-certificates

# STEP 9: Set working directory
WORKDIR /app

# STEP 10: Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# STEP 11: Copy .env file for environment variables
COPY .env .

# STEP 12: Expose the required port
EXPOSE 8080

# STEP 13: Command to run the application
CMD ["./main"]

# Use the official Go image as the base image
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install Go dependencies
RUN go mod download

# Copy the entire application directory to the working directory
COPY . .

# Build the Go application with support for embedding
RUN go build -o app

# Create a smaller image without the source code and unnecessary build artifacts
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the binary from the build image
COPY --from=build /app/app .

# Expose the port that your application will run on
EXPOSE 8080

# Command to run the application
CMD ["./app"]
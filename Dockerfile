FROM golang:latest

# Add Maintainer Info
LABEL maintainer="container for local test"
# Set the Current Working Directory inside the container
WORKDIR /mafiosi
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . ./
# Build the Go app and Command to run the executable
RUN go build -o app .
EXPOSE 8080
# Use -host flag if needed "./app -host=yourhost:port" default:"localhost:8080"
CMD ["./app"]

# Start from golang:1.12-alpine base image
FROM golang:1.12-alpine

LABEL maintainer="bit-framework"
# Set the Current Working Directory inside the container
WORKDIR /bit-framework

# Copy go mod and sum files
COPY cmd/go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd .

# Build the Go app
RUN go build -o main .

# Expose port to the outside world
EXPOSE 8079
EXPOSE 8081

# Run the executable
CMD ["./cmd"]
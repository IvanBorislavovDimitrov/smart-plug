FROM golang:alpine

LABEL maintainer="IvanBorislavovDimitrov"

ENV CONN_STR "postgresql://postgres:123@localhost:5433/smart_plug?sslmode=disable"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Install Curl
RUN apk update && apk add --no-cache curl

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY ./ /app/

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /build

# Expose port 8081 to the outside world
EXPOSE 8081

# Run the executable
CMD [ "./build.sh" ]
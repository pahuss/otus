#FROM ubuntu:latest
#
#ARG GO_VERSION=1.16.7
#ENV GO_VERSION=${GO_VERSION}
#
#RUN apt-get update
#RUN apt-get install -y wget git gcc
#
#RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"
#
#RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
#RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
#
#ENV GOPATH /go
#ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
#RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
#
#WORKDIR $GOPATH

# Dockerfile References: https://docs.docker.com/engine/reference/builder/
#################################################################################################################
## Start from golang:1.12-alpine base image
FROM golang:alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Add Maintainer Info
LABEL maintainer="G P <pahuss@mail.ru>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
#RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
#CMD ["./main"]
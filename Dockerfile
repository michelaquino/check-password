# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest

# Copy the local package files to the container's workspace.
ADD . /go/src/gitlab.globoi.com/michel.aquino/check-password

# Build the api command inside the container.
WORKDIR /go/src/gitlab.globoi.com/michel.aquino/check-password
RUN go build -o app

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/src/gitlab.globoi.com/michel.aquino/check-password/app

# Document that the service listens on port 8888.
EXPOSE 8888
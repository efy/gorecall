FROM golang:1.9.2

# Ensure go bin directory exists
RUN mkdir -p /go/bin

# Make and set the working directory
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

# Add app source to working directory
ADD . /go/src/app

# Fetch dependencies
RUN go get -u github.com/golang/dep/...
RUN dep ensure

# Build the app
RUN go build -o /go/bin/recall /go/src/app/main.go

# Creates a docker image for building linux binaries
# using go dep for dependency managemenet

FROM golang:1.9.2

# Install go dep
ADD https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 ./bin/dep

# Make dep executable
RUN chmod +x ./bin/dep



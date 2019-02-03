# Golang build
FROM golang:1.11 as server_builder
COPY . /usr/src/go
WORKDIR /usr/src/go

RUN GOOS=linux GOARCH=amd64 go build -o recall main.go

# Web client build env
FROM node:9.6.1 as client_builder
RUN mkdir /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH
COPY ./client/package.json /usr/src/app/package.json
RUN npm install --silent
RUN npm install react-scripts@1.1.1 -g --silent
COPY ./client /usr/src/app
RUN npm run build

# Prod env
FROM alpine:3.7

# Add certificates
RUN apk add --no-cache ca-certificates

COPY --from=client_builder /usr/src/app/build /usr/share/recall-client
COPY --from=server_builder /usr/src/go/recall /usr/bin/recall
EXPOSE 8080
CMD ["/usr/bin/recall", "serve"]

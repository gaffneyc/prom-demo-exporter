# Intermediary container to build the application's container
FROM golang:alpine AS build

# Required for CGO support
RUN apk --no-cache add gcc musl-dev git

COPY . /prom-demo-exporter/
WORKDIR /prom-demo-exporter
RUN GOOS=linux GOARCH=amd64 go build ./...

# Build the final image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

EXPOSE 1845

COPY --from=build /prom-demo-exporter/prom-demo-exporter /bin

CMD [ "/bin/prom-demo-exporter" ]

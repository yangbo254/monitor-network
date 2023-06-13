FROM golang:1.18 AS build-stage

WORKDIR /usr/local/go/src/nginxbackend

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /monitor-network

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...


# Deploy the application binary into a lean image
FROM debian:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /monitor-network /monitor-network
COPY ["html","html"]

ENTRYPOINT ["/monitor-network"]
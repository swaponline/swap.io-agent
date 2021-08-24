FROM golang:1.16-alpine AS build
WORKDIR /go/src/swap.io-agent
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /build

FROM ubuntu
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build .
USER root:root
CMD [ "./build" ]
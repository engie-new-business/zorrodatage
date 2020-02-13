FROM golang:1.13.4 as builder

RUN mkdir -p build

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -trimpath -ldflags '-s -w'

FROM ubuntu:18.04
# update certificates for external call to HTTPS service
RUN apt-get update && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/zorrodatage .

ENTRYPOINT ["./zorrodatage"]
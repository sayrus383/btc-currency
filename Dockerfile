#--------------------------------------
FROM golang:1.21.5-alpine3.17 as cache

RUN apk add --update --no-cache \
    tzdata \
    ca-certificates

RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

#--------------------------------------
FROM cache as build

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -v -ldflags="-s -w" -o /usr/bin/btc-currency /app/cmd/main.go

#--------------------------------------
FROM alpine:3.17 as app

COPY --from=build /usr/bin/btc-currency /usr/bin/btc-currency
RUN chmod +x /usr/bin/btc-currency

CMD ["/usr/bin/btc-currency"]

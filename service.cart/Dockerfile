FROM golang:1.16.3-alpine3.13 as builder


RUN go env -w GO111MODULE=on

# For china proxy
RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o service-cart /build/service.cart/cmd/main.go

# FROM scratch
FROM alpine:3.13.4

WORKDIR /app

# add timezone
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

COPY --from=builder /build/config /app/config
COPY --from=builder /build/service-cart /app/service-cart
COPY --from=builder /build/migrations /app/migrations

EXPOSE 5000

CMD ["./service-cart"]

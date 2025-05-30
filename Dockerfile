# Этап сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prometheus-temperature-disk .

FROM alpine:3.21

ENV TZ=Europe/Moscow
RUN apk add --update --no-cache tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

RUN apk add --no-cache smartmontools
COPY --from=builder /app/prometheus-temperature-disk /app/prometheus-temperature-disk
EXPOSE 9586
ENTRYPOINT ["/app/prometheus-temperature-disk"]

# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prometheus-disk-temp .

FROM alpine:3.21

ENV TZ=Europe/Moscow
RUN apk add --update --no-cache tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

RUN apk add --no-cache smartmontools nvme-cli
COPY --from=builder /app/prometheus-disk-temp /app/prometheus-disk-temp
EXPOSE 9586
CMD ["/app/prometheus-disk-temp"]

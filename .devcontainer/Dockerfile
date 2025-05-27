FROM alpine:3.21

ENV TZ=Europe/Moscow
RUN apk add --update --no-cache tzdata \
    && cp /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

RUN apk add --update --no-cache go
RUN apk add --update --no-cache smartmontools

# Set `DEVCONTAINER` environment variable to help with orientation
RUN apk add --update --no-cache git
ENV DEVCONTAINER=true
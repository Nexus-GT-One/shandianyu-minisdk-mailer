FROM alpine:3.22
WORKDIR /srv/app
RUN apk update && apk upgrade --available && apk add --no-cache -U tzdata busybox-extras && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone
COPY ./config.yaml ./config.yaml
COPY ./shandianyu-minisdk-mailer ./shandianyu-minisdk-mailer
CMD ["./shandianyu-minisdk-mailer"]
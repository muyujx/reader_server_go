FROM alpine:3.22

WORKDIR /reader
COPY dist/reader /reader
COPY config.yml /reader

RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && chmod +x ./reader

EXPOSE 8080
VOLUME ["/data"]

ENTRYPOINT ["./reader"]

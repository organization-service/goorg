FROM golang:1.14-alpine3.11 as build
RUN apk --no-cache add tzdata gcc libc-dev git openssh \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && apk del tzdata \
    && go get -v github.com/rubenv/sql-migrate/... \
    && rm  -rf /tmp/* /var/cache/apk/*
WORKDIR /var/app/golang
COPY backend .
RUN go get -v \
    && go build -o app

FROM alpine:3
ENV TZ=Asia/Tokyo
COPY --from=build /var/app/golang/app /app/
COPY --from=build /go/bin/sql-migrate /usr/bin/
COPY backend/dbconfig.yml /app/
COPY backend/migrations/ /app/migrations/
COPY docker-entrypoint.sh /usr/bin/docker-entrypoint.sh
WORKDIR /app
RUN chmod +x /app/app \
    && chmod +x /usr/bin/docker-entrypoint.sh \
    && chmod +x /usr/bin/sql-migrate \
    && apk --no-cache add tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && rm  -rf /tmp/* /var/cache/apk/*

EXPOSE 8080

ENTRYPOINT [ "docker-entrypoint.sh" ]

CMD [ "./app" ]
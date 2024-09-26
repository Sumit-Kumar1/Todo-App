FROM alpine:latest

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite

WORKDIR /todoApp

RUN mkdir -p {views,public}

COPY views/*.html ./views/
COPY public/*  ./public/
COPY .env .
COPY Build/main .

RUN chmod +x main

RUN chown -R nonroot:nonroot .
USER nonroot

ENTRYPOINT [ "./main" ]

EXPOSE 9001
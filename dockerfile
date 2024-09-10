FROM alpine:latest

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

RUN apk update && apk upgrade
RUN apk add --no-cache sqlite

WORKDIR /todoApp

RUN mkdir views

COPY views/index.html ./views/
COPY Build/main .

RUN chmod +x main

RUN touch tasks.db

RUN chown -R nonroot:nonroot .

USER nonroot

ENTRYPOINT [ "./main" ]

EXPOSE 9001
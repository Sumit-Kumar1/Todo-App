FROM alpine:latest

WORKDIR /todoApp

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot && \
apk update && apk upgrade

COPY views/*.html ./views/
COPY public/*  ./public/
COPY .env ./
COPY Build/main .

RUN chmod +x main && chown -R nonroot:nonroot .
USER nonroot

ENTRYPOINT [ "./main" ]

EXPOSE 9001
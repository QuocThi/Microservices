FROM alpine:latest

RUN mkdir /app

COPY randomService /app

CMD [ "/app/randomService"]

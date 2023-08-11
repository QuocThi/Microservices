FROM alpine:latest

RUN mkdir /app

COPY randomApp /app

CMD [ "/app/randomApp"]

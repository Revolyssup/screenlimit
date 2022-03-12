FROM ubuntu:latest

COPY . /app

WORKDIR /app

ENTRYPOINT [ "./screenlim" ]
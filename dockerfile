FROM golang:latest

RUN apt-get update
RUN apt-get install -y chromium

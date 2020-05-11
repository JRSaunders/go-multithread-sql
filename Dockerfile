FROM ubuntu:18.04

WORKDIR /

COPY . /
ENV GOTHREADED_DOCKER=true
CMD ["./gothreaded"]

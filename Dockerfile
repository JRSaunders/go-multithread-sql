FROM amd64/ubuntu:20.04

WORKDIR /

COPY . /
ENV GOTHREADED_DOCKER=true
CMD ["./gothreaded"]

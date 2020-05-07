FROM ubuntu:18.04

WORKDIR /

COPY . /

CMD ["./gothreaded"]

FROM golang:1.17.9-stretch as builder

RUN mkdir -p /go/instancer

COPY * /go/instancer/
WORKDIR /go/instancer/

RUN go build .

EXPOSE 8888/tcp

ENV DOCKER_HOST "unix:///var/run/docker.sock"

CMD [ "/go/instancer/minetest_instancer", "8888"]

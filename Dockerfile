FROM golang:1.16.3
WORKDIR /dcard-demo
ADD . /dcard-demo
RUN cd /dcard-demo && go build
ENTRYPOINT ./dcard-demo
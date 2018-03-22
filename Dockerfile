FROM golang:alpine as builder

ENV PATH=${PATH}:${GOPATH}/bin

RUN apk update && apk add git make

COPY . ${GOPATH}/src/github.com/while-loop/todo
WORKDIR ${GOPATH}/src/github.com/while-loop/todo
RUN make all

FROM alpine:latest
COPY --from=builder /go/bin/todod /usr/local/bin/
RUN todod -v

CMD ["todod"]
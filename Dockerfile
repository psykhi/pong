FROM golang:alpine as builder
ARG package_root=github.com/psykhi/pong
WORKDIR /go/src/$package_root
COPY . /go/src/$package_root
RUN go build  -o /go/bin/server /go/src/$package_root/server/cmd/

FROM alpine
COPY --from=builder /go/bin/server /
CMD ["/server"]

FROM alpine:3.15
ADD ./tkdo /go/bin/
ENTRYPOINT /go/bin/tkdo

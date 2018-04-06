FROM golang:alpine as builder
LABEL maintainer="Ben Selby <benmatselby@gmail.com>"

# RUN	apk add --no-cache \
# 	ca-certificates

COPY . /go/src/github.com/benmatselby/donny

RUN apk update && \
    apk add gcc libc-dev libgcc git make

RUN cd /go/src/github.com/benmatselby/donny \
	&& make static \
	&& mv donny /usr/bin/donny \
	&& rm -rf /go

FROM scratch

COPY --from=builder /usr/bin/donny /usr/bin/donny

ENTRYPOINT [ "donny" ]

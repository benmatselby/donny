FROM golang:alpine as builder
LABEL maintainer="Ben Selby <benmatselby@gmail.com>"

ENV APPNAME donny
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go
ENV DEP_VERSION 0.4.1

COPY . /go/src/github.com/benmatselby/${APPNAME}

RUN apk update && \
    apk add --no-cache --virtual .build-deps \
		ca-certificates \
		gcc \
		libc-dev \
		libgcc \
		git \
		curl \
		make

RUN curl -fsSL -o /usr/local/bin/dep \
	https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 && \
	chmod +x /usr/local/bin/dep

RUN cd /go/src/github.com/benmatselby/${APPNAME} && \
	make static-all  && \
	mv ${APPNAME} /usr/bin/${APPNAME}  && \
	apk del .build-deps  && \
	rm -rf /go

FROM scratch

COPY --from=builder /usr/bin/${APPNAME} /usr/bin/${APPNAME}
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "donny" ]
CMD [ "--help"]

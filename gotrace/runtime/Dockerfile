FROM buildpack-deps:jessie-scm

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		patch \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.4.3
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 ce3140662f45356eb78bc16a88fc7cfb29fb00e18d7c632608245b789b2086d2

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/gotemp/go/bin:$PATH

ENV GOLANG_TO_BUILD_VERSION 1.8.1
ENV GOLANG_TO_BUILD_URL https://golang.org/dl/go$GOLANG_TO_BUILD_VERSION.src.tar.gz
ENV GOLANG_TO_BUILD_SHA256 33daf4c03f86120fdfdc66bddf6bfff4661c7ca11c5da473e537f4d69b470e57

ENV GOROOT_BOOTSTRAP /usr/local/gotemp/go

ADD . /runtime

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& mkdir -p /usr/local/gotemp && tar -C /usr/local/gotemp -xzf golang.tar.gz \
	&& rm golang.tar.gz \
	&& mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" \
	&& curl -fsSL "$GOLANG_TO_BUILD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_TO_BUILD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz \
	&& cd /usr/local/go/src \
	&& patch -p2 < /runtime/go"$GOLANG_TO_BUILD_VERSION"-gotrace.patch \
	&& ./make.bash \
	&& rm -rf /usr/local/gotemp/

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH

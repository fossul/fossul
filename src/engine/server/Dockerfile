FROM fedora:latest


LABEL ios.k8s.display-name="fossul-server" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir -p /go/src/github.com/fossul/fossul
RUN mkdir /app

RUN groupadd -r fossul &&\
    useradd -r -g fossul -d /app -u 1001 -s /sbin/nologin -c "Fossul user" fossul

WORKDIR /go/src/github.com/fossul/fossul

ENV FOSSUL_BUILD_PLUGIN_DIR=/app/plugins
ENV GOPATH=/go
ENV GOBIN=/app
ENV HOME=/app

RUN curl -L https://github.com/fossul/fossul/releases/download/latest/openshift-client-linux-4.2.7.tar.gz |tar xz;cp oc kubectl /app && \

RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh

COPY . /go/src/github.com/fossul/fossul

RUN /go/src/github.com/fossul/fossul/scripts/fossul-server-build.sh

RUN chown -R fossul:fossul /app && \
    chmod -R 775 /app

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER fossul

WORKDIR /app

CMD /app/fossul-server-startup.sh

FROM fedora:latest


LABEL ios.k8s.display-name="fossul-storage" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir -p /go/src/fossul
RUN mkdir /app
RUN mkdir /app/backups

RUN groupadd -r fossul &&\
    useradd -r -g fossul -d /app -u 1001 -s /sbin/nologin -c "Fossul user" fossul &&\
    echo 'fossul ALL=(ALL) NOPASSWD: ALL' >/etc/sudoers.d/fossul

WORKDIR /go/src/fossul

ENV FOSSUL_BUILD_PLUGIN_DIR=/app/plugins
ENV GOPATH=/go
ENV GOBIN=/app
ENV HOME=/app

RUN curl -L https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v3.11.0-0cbc58b-linux-64bit.tar.gz |tar xz;cp openshift*/oc /app;cp openshift*/kubectl /app

RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh

COPY . /go/src/fossul

RUN /go/src/fossul/scripts/fossul-storage-build.sh

RUN chown -R fossul:fossul /app && \
    chmod -R 775 /app

RUN chmod -R 777 /tmp

RUN ls -l /app

RUN echo "1.0" > /etc/imageversion

USER fossul

WORKDIR /app

CMD /app/fossul-storage-startup.sh

FROM fedora:latest


LABEL ios.k8s.display-name="fossil-storage" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir -p /go/src/fossil
RUN mkdir /app
RUN mkdir /app/backups

WORKDIR /go/src/fossil

ENV FOSSIL_BUILD_PLUGIN_DIR=/app/plugins
ENV GOPATH=/go
ENV GOBIN=/app

RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh

#RUN go get -u github.com/gorilla/mux

COPY . /go/src/fossil

RUN /go/src/fossil/fossil-storage-build.sh

#RUN cd /app && go build engine/util && go install engine/storage

RUN chown -R 1001:0 /app && \
    chmod -R 777 /app

RUN chmod -R 777 /tmp

RUN ls -l /app

RUN echo "1.0" > /etc/imageversion

USER 1001

WORKDIR /app

CMD /app/fossil-storage-startup.sh

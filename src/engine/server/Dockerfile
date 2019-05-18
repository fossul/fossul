FROM fedora:latest


LABEL ios.k8s.display-name="fossil-server" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir -p /go/src/fossil
RUN mkdir /app
RUN mkdir /app/data
RUN mkdir -p /app/configs/default

WORKDIR /go/src/fossil

ENV FOSSIL_BUILD_PLUGIN_DIR=/app/plugins
ENV GOPATH=/go
ENV GOBIN=/app

RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh

#RUN go get -u github.com/gorilla/mux

COPY . /go/src/fossil

RUN /go/src/fossil/fossil-server-build.sh

#RUN cd /app && go build engine/util && go install engine/workflow

RUN chown -R 1001:0 /app && \
    chmod -R 777 /app

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /app/fossil-server-startup.sh

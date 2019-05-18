FROM fedora:latest


LABEL ios.k8s.display-name="fossil-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir -p /go/src/fossil
RUN mkdir /app

WORKDIR /app

ENV FOSSIL_BUILD_PLUGIN_DIR=/app/plugins
ENV GOPATH=/go
ENV GOBIN=/app

RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh

#RUN go get -u github.com/gorilla/mux

COPY . /go/src/fossil

#RUN cd /app && go build engine/util && go install engine/app
RUN /go/src/fossil/fossil-app-build.sh

RUN chown -R 1001:0 /app && \
    chmod -R 777 /app
 #   chmod -R ug+rwX /app

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

RUN ls -l /app

USER 1001

CMD /app/fossil-app-startup.sh

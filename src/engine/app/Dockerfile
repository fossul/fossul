FROM fedora:latest


LABEL ios.k8s.display-name="fossil-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir /app
RUN mkdir /app/go
RUN mkdir /app/go/src
RUN mkdir /app/go/src/fossil

WORKDIR /app

ENV FOSSIL_BUILD_PLUGIN_DIR=/app
ENV GOPATH=/app/go
ENV GOBIN=/app

RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh

#RUN go get -u github.com/gorilla/mux

COPY . /app/go/src/fossil

#RUN cd /app && go build engine/util && go install engine/app
RUN /bin/bash /app/go/src/fossil/fossil-app-build.sh

RUN chown -R 1001:0 /app && \
    chmod -R ug+rwX /app && \
    chmod 755 /app/app

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /bin/bash /app/fossil-app-startup.sh

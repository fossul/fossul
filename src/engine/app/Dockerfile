FROM fedora:latest


LABEL ios.k8s.display-name="fossil-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir -p /app/src

WORKDIR /app

ENV GOPATH=/app
ENV GOBIN=/app

RUN go get -u github.com/gorilla/mux

COPY . /app

#RUN cd /app/engine/app && go build -o /app/app main.go logger.go routes.go handlers.go structs.go
RUN cd /app && go build src/engine/util && go install src/engine/app

RUN chown -R 1001:0 /app && \
    chmod -R ug+rwX /app && \
    chmod 755 /app/app

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /app/app

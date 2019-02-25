FROM fedora:latest


LABEL ios.k8s.display-name="fossil-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir /app

WORKDIR /app

ENV GOPATH=/app

RUN go get -u github.com/gorilla/mux

COPY . /app

RUN cd /app/engine/workflow && go build -o /app/workflow main.go logger.go routes.go handlers.go structs.go

RUN chown -R 1001:0 /app && \
    chmod -R ug+rwX /app 

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /bin/bash /app/workflow

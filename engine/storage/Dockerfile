FROM fedora:latest


LABEL ios.k8s.display-name="fossil-storage" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN go get -u github.com/gorilla/mux

RUN mkdir /backups

RUN mkdir /app

WORKDIR /app

COPY . /app

RUN cd /app/engine/storage && go build -o /app/storage main.go logger.go routes.go handlers.go structs.go

RUN chown -R 1001:0 /app && \
    chown -R 1001:0 /backups && \
    chmod -R ug+rwX /app &&

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /bin/bash /app/storage
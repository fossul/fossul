FROM fedora:latest


LABEL ios.k8s.display-name="fossil-storage" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install golang

RUN go build main.go logger.go routes.go handlers.go struct.go -o storage

RUN mkdir /backups

RUN mkdir /app

WORKDIR /app

COPY . /app


RUN chown -R 1001:0 /app && \
    chown -R 1001:0 /backups && \
    chmod -R ug+rwX /app

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /bin/bash /app/storage
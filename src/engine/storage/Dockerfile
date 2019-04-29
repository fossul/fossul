FROM fedora:latest


LABEL ios.k8s.display-name="fossil-storage" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN dnf install -y golang

RUN mkdir /app

WORKDIR /app

ENV GOPATH=/app
ENV GOBIN=/app

#RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh

RUN go get -u github.com/gorilla/mux

RUN mkdir /backups

COPY . /app

RUN cd /app && go build engine/util && go install engine/storage

RUN chown -R 1001:0 /app && \
    chown -R 1001:0 /backups && \
    chmod -R ug+rwX /app && \
    chmod 755 /app/storage

RUN chmod -R 777 /tmp

RUN echo "1.0" > /etc/imageversion

USER 1001

CMD /app/storage

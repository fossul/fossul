FROM registry.access.redhat.com/ubi8/go-toolset:latest

LABEL ios.k8s.display-name="fossul-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

RUN mkdir /opt/app-root/fossul
ENV GOBIN=/opt/app-root/fossul
ENV APP_PLUGIN_DIR=/opt/app-root/fossul/plugins/app

#RUN curl -L https://github.com/fossul/fossul/releases/download/latest/openshift-client-linux-4.2.8.tar.gz |tar xz;cp oc kubectl /app

COPY . /opt/app-root/src/github.com/fossul/fossul

WORKDIR /opt/app-root/src/github.com/fossul/fossul
RUN /opt/app-root/src/github.com/fossul/fossul/scripts/fossul-app-build.sh

USER 1001
WORKDIR /opt/app-root/fossul
CMD ./fossul-app-startup.sh

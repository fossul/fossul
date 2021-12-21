FROM registry.access.redhat.com/ubi8/go-toolset:latest

LABEL ios.k8s.display-name="fossul-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root
ENV STORAGE_PLUGIN_DIR=/opt/app-root/plugins/storage
ENV ARCHIVE_PLUGIN_DIR=/opt/app-root/plugins/archive

#RUN curl -L https://github.com/fossul/fossul/releases/download/latest/openshift-client-linux-4.2.8.tar.gz |tar xz;cp oc kubectl /app

COPY . /opt/app-root/src/github.com/fossul/fossul

WORKDIR /opt/app-root/src/github.com/fossul/fossul
RUN /opt/app-root/src/github.com/fossul/fossul/scripts/fossul-storage-build.sh

USER 1001
WORKDIR /opt/app-root
CMD /opt/app-root/fossul-storage-startup.sh

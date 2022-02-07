FROM registry.access.redhat.com/ubi8/go-toolset:latest
LABEL ios.k8s.display-name="fossul-server-build" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root/src
ENV APP_PLUGIN_DIR=/opt/app-root/fossul/plugins/app
#RUN curl -L https://github.com/fossul/fossul/releases/download/latest/openshift-client-linux-4.2.8.tar.gz |tar xz;cp oc kubectl /app
COPY . /opt/app-root/src/github.com/fossul/fossul
WORKDIR /opt/app-root/src/github.com/fossul/fossul
RUN /opt/app-root/src/github.com/fossul/fossul/scripts/fossul-server-build.sh
#RUN mkdir -p /opt/app-root/metadata/configs/default/default
#COPY src/cli/configs/default /opt/app-root/metadata/configs/default/default

FROM registry.access.redhat.com/ubi8/ubi:latest
LABEL ios.k8s.display-name="fossul-server" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root
RUN mkdir -p /opt/app-root
WORKDIR /opt/app-root
COPY --from=0 /opt/app-root/src/server ./
COPY --from=0 /opt/app-root/src/fossul-server-startup.sh ./
RUN mkdir -p metadata/configs
RUN mkdir -p metadata/data
RUN chown -R 1001:0 /opt/app-root

USER 1001
CMD ./fossul-server-startup.sh

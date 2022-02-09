FROM registry.access.redhat.com/ubi8/go-toolset:latest
LABEL ios.k8s.display-name="fossul-server-build" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root/src
ENV APP_PLUGIN_DIR=/opt/app-root/fossul/plugins/app

COPY . /opt/app-root/src/github.com/fossul/fossul
WORKDIR /opt/app-root/src/github.com/fossul/fossul
RUN /opt/app-root/src/github.com/fossul/fossul/scripts/fossul-server-build.sh


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
RUN mkdir -p /opt/app-root/default
COPY src/cli/configs/default /opt/app-root/default
RUN chown -R 1001:0 /opt/app-root

USER 1001
CMD ./fossul-server-startup.sh

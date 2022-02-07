FROM registry.access.redhat.com/ubi8/go-toolset:latest

LABEL ios.k8s.display-name="fossul-app-build" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root/src
ENV APP_PLUGIN_DIR=/opt/app-root/src/plugins/app
#RUN curl -L https://github.com/fossul/fossul/releases/download/latest/openshift-client-linux-4.2.8.tar.gz |tar xz;cp oc kubectl /app
COPY . /opt/app-root/src/github.com/fossul/fossul
WORKDIR /opt/app-root/src/github.com/fossul/fossul
RUN /opt/app-root/src/github.com/fossul/fossul/scripts/fossul-app-build.sh

FROM registry.access.redhat.com/ubi8/ubi:latest
LABEL ios.k8s.display-name="fossul-app" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root
ENV APP_PLUGIN_DIR=/opt/app-root/plugins/app
RUN mkdir -p /opt/app-root/plugins
WORKDIR /opt/app-root
COPY --from=0 /opt/app-root/src/app ./
COPY --from=0 /opt/app-root/src/fossul-app-startup.sh ./
COPY --from=0 /opt/app-root/src/plugins ./plugins
RUN chown -R 1001:0 /opt/app-root


USER 1001
CMD ./fossul-app-startup.sh

FROM registry.access.redhat.com/ubi8/go-toolset:latest
LABEL ios.k8s.display-name="fossul-cli-build" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root/src
COPY . /opt/app-root/src/github.com/fossul/fossul
WORKDIR /opt/app-root/src/github.com/fossul/fossul
RUN /opt/app-root/src/github.com/fossul/fossul/scripts/fossul-cli-build.sh

FROM registry.access.redhat.com/ubi8/ubi:latest
LABEL ios.k8s.display-name="fossul-cli" \
    maintainer="Keith Tenzer <ktenzer@redhat.com>"

ENV GOBIN=/opt/app-root
ENV HOME=/opt/app-root
RUN mkdir -p /opt/app-root
WORKDIR /opt/app-root
COPY --from=0 /opt/app-root/src/cli ./
COPY --from=0 /opt/app-root/src/fossul-cli-startup.sh ./
RUN chown -R 1001:0 /opt/app-root

USER 1001
CMD ./fossul-cli-startup.sh

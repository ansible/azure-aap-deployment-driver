FROM registry.access.redhat.com/ubi9/nginx-120

USER root

WORKDIR /opt/app-root/src/

ARG ACME_RELEASE_TAG=3.0.5
ARG DRIVER_RELEASE_TAG=dev

# Install socat tools required by acme.sh and download acme.sh from version/tag (but don't install)
RUN yum -y --repo ubi-9-appstream-rpms install socat && \
  curl -L -o acme.zip https://github.com/acmesh-official/acme.sh/archive/refs/tags/v${ACME_RELEASE_TAG}.zip && \
  unzip -qoj acme.zip acme.sh-${ACME_RELEASE_TAG}/acme.sh -d . && rm acme.zip && \
  echo "ACME=${ACME_RELEASE_TAG}" >> versions && echo "DRIVER=${DRIVER_RELEASE_TAG}" >> versions && ls -al && cat versions

ADD ["nginx", "/etc/nginx/"]
ADD ["start.sh", "build/server", "./"]
ADD ["build/public", "/var/www/aapinstaller/public"]

RUN chmod +x ./acme.sh ./server && chmod +x ./start.sh

VOLUME [ "/installerstore" ]

ENTRYPOINT ["./start.sh"]

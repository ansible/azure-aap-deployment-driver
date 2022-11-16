FROM registry.access.redhat.com/ubi9/nginx-120

USER root

WORKDIR /opt/app-root/src/

# Install socat tools required by acme.sh and download acme.sh (but don't install)
RUN yum -y --repo ubi-9-appstream-rpms install socat && \
  curl -o acme.installer.sh https://get.acme.sh && chmod +x acme.installer.sh

ADD ["nginx", "/etc/nginx/"]
ADD ["start.sh", "build/server", "./"]
ADD ["build/templates/", "./templates/"]
ADD ["build/public", "/var/www/aapinstaller/public"]

RUN chmod +x ./server && chmod +x ./start.sh

VOLUME [ "/installerstore" ]

ENTRYPOINT ["./start.sh"]

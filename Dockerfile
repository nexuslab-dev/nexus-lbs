FROM ubuntu:22.04

RUN mkdir /etc/nexus-lbs; \
    apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates tzdata; \
    update-ca-certificates -f; \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \
    apt autoremove -y; \
    rm -rf /var/lib/apt/lists/*

COPY nexus-lbs  /usr/local/bin/

COPY ipdb/GeoLite2-City.mmdb /etc/nexus-lbs/

COPY ipdb/GeoLite2-Country.mmdb /etc/nexus-lbs/

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

WORKDIR /usr/local/bin/

ENV TZ=Asia/Shanghai \
SENTRY_DSN="" \
SENTRY_ENVIRONMENT="" \
SENTRYGODEBUG=""

VOLUME /etc/nexus-lbs

ENTRYPOINT ["/tini", "--"]

# Run your program under Tini
CMD ["/usr/local/bin/nexus-lbs", "-c", "/etc/nexus-lbs/config.toml"]

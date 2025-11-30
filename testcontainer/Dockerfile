FROM debian:bookworm-slim
ARG DEBIAN_FRONTEND=noninteractive

# Locations
ARG INSTALL_PATH=/opt/snitchdns
ARG DATA_PATH=$INSTALL_PATH/data
ARG CONFIG_PATH=$DATA_PATH/config
ARG REPO=https://github.com/sadreck/SnitchDNS

# Configuration
ARG SNITCHDNS_DBMS=sqlite
ARG SNITCHDNS_SECRET_KEY=TestSecretKeyForContainerOnlyNotForProductionUseThisMustBeLongEnough
ARG SNITCH_DOMAIN=snitch.test
ARG BASE_DOMAIN=snitch.test

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8

# Install prerequisites
RUN apt-get update && \
    apt-get install -y \
    git \
    python3-pip \
    python3-venv \
    libpq-dev \
    cron \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install SnitchDNS
RUN git clone $REPO /opt/snitchdns

WORKDIR $INSTALL_PATH

RUN python3 -m venv venv && \
    . venv/bin/activate && \
    pip --no-cache-dir install 'Flask-Migrate<3.0' && \
    pip --no-cache-dir install -r requirements.txt && \
    deactivate

# Configure SnitchDNS
RUN mkdir -p $CONFIG_PATH/env && \
    mkdir -p $CONFIG_PATH/service && \
    mkdir -p $CONFIG_PATH/http

RUN touch $CONFIG_PATH/env/snitch.conf && \
    echo SNITCHDNS_DBMS=$SNITCHDNS_DBMS > $CONFIG_PATH/env/snitch.conf && \
    echo SNITCHDNS_DB_USER=none >> $CONFIG_PATH/env/snitch.conf && \
    echo SNITCHDNS_DB_PW=none >> $CONFIG_PATH/env/snitch.conf && \
    echo SNITCHDNS_DB_URL=none >> $CONFIG_PATH/env/snitch.conf && \
    echo SNITCHDNS_DB_DB=none >> $CONFIG_PATH/env/snitch.conf && \
    echo SNITCHDNS_SECRET_KEY=$SNITCHDNS_SECRET_KEY >> $CONFIG_PATH/env/snitch.conf

RUN ln -s $CONFIG_PATH/env/snitch.conf $INSTALL_PATH/.env

# Initialize database
RUN ./venv.sh flask db init && \
    ./venv.sh flask db migrate && \
    ./venv.sh flask db upgrade && \
    ./venv.sh flask snitchdb && \
    ./venv.sh flask crontab add

# Configure DNS daemon
RUN ./venv.sh flask settings set --name dns_daemon_bind_ip --value 0.0.0.0 && \
    ./venv.sh flask settings set --name dns_daemon_bind_port --value 2024 && \
    ./venv.sh flask settings set --name dns_base_domain --value $BASE_DOMAIN

# Create test user and API key via entrypoint
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

VOLUME $DATA_PATH

EXPOSE 80 2024/udp

HEALTHCHECK --interval=10s --timeout=5s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:80/ || exit 1

ENTRYPOINT ["/entrypoint.sh"]

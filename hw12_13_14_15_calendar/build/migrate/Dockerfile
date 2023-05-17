FROM ubuntu:20.04

ENV PG_URL='postgres://postgres:password@localhost:5432/example?sslmode=disable'
ENV DIR=/migrations

ARG APP_USERNAME=migrator
ARG APP_USER_UID=1000

RUN apt-get -y update \
  && apt-get -y install curl add-apt-key lsb-release \
  && curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add - \
  && echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list \
  && apt-get update && apt-get install -y migrate && apt-get clean all \
  && useradd --create-home --uid ${APP_USER_UID} ${APP_USERNAME} && mkdir -p ${DIR} && chown -R ${APP_USERNAME}:${APP_USERNAME} ${DIR}

COPY --chown=${APP_USERNAME}:${APP_USERNAME} ./migrations ${DIR}

USER ${APP_USERNAME}

CMD migrate -path ${DIR} -database ${PG_URL} -verbose up

FROM ubuntu:18.04

EXPOSE 5432
EXPOSE 5000

ENV DEBIAN_FRONTEND 'noninteractive'
ENV PGVER 10

RUN apt -y update && apt install -y postgresql-$PGVER
RUN apt install -y wget
RUN apt install -y git

USER postgres

ADD ./init.sql /opt/init.sql

RUN /etc/init.d/postgresql start &&\
	psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
	createdb -O docker docker &&\
    psql -f /opt/init.sql -d docker &&\
    /etc/init.d/postgresql stop

RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "include_dir='conf.d'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/$PGVER/main/postgresql.conf

ADD ./postgres.conf /etc/posуtgresql/$PGVER/main/conf.d/basic.conf

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

RUN wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
RUN tar -xvf go1.13.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

ADD . /opt/app
WORKDIR /opt/app
CMD service postgresql start && go run main.go

FROM ubuntu

RUN apt-get -y update

RUN apt-get -y install npm

RUN npx playwright install --with-deps

COPY cacert.pem /etc/ssl/certs/

COPY config.ini /etc/project/
COPY main /

CMD ["/main"]
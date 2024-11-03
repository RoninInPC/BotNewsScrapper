FROM ubuntu

RUN apt-get -y update

RUN apt-get -y install npm
RUN npm install
RUN npm install -g playwright
RUN npx -y playwright install --with-deps

COPY cacert.pem /etc/ssl/certs/

COPY config.ini /etc/project/
COPY main /

CMD ["/main"]
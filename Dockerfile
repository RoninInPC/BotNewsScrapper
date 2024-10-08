FROM ubuntu

RUN apt-get -y update

RUN apt-get -y install npm

RUN npx playwright install-deps

RUN apt-get install libatk1.0-0t64  \
  libatk-bridge2.0-0t64  \
    libcups2t64  \
    libatspi2.0-0t64  \
    libxcomposite1  \
    libxdamage1  \
    libxfixes3 \
    libxrandr2 \
    libgbm1 \
    libpango-1.0-0 \
    libcairo2 \
    libasound2t64

COPY cacert.pem /etc/ssl/certs/

COPY config.ini /etc/project/
COPY main /

CMD ["/main"]
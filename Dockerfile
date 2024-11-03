FROM node:20-bookworm

RUN npx -y playwright@1.47.2 install --with-deps

COPY cacert.pem /etc/ssl/certs/

COPY config.ini /etc/project/
COPY main /

CMD ["/main"]
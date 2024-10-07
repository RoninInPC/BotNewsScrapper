FROM ubuntu

COPY cacert.pem /etc/ssl/certs/

COPY config.ini /
COPY main /

CMD ["/main"]

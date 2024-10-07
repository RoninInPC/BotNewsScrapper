FROM ubuntu


COPY cacert.pem /etc/ssl/certs/

COPY config.ini /etc/project/
COPY main /

RUN ./main

CMD ["/main"]
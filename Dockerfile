FROM ubuntu


COPY cacert.pem /etc/ssl/certs/

COPY config.ini /etc/project/
COPY main /

RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

RUN playwright install --with-deps

CMD ["/main"]

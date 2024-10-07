FROM ubuntu

RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

RUN playwright install --with-deps

COPY cacert.pem /etc/ssl/certs/

COPY config.ini /etc/project/
COPY main /

CMD ["/main"]

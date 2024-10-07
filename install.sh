#!/bin/bash

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
export PATH=$PATH:/usr/local/go/bin &&
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . &&
docker build -t bot_docker .

cd service &&
cp bot.service /etc/systemd/system/ &&
cp redis.service /etc/systemd/system/ &&
cd /etc/systemd/system &&
systemctl daemon-reload &&
systemctl restart redis.service && systemctl restart bot.service &&
systemctl enable redis.service && systemctl enable bot.service



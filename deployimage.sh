#!/bin/bash
echo "inserted user:$1 host:$2 and port:$3 "

cd makefiles; make build;
docker save yt-telegram-bot:latest | gzip > yt-telegram-bot.tar.gz
scp -P $3 ./yt-telegram-bot.tar.gz $1@$2:/$1 
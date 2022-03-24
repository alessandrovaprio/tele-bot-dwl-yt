- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Run Docker](#run-docker)
  - [Run Binary](#run-binary)
- [Libraries](#libraries)
- [License](#license)

# Overview
A bot useful to download Video or Audio From Youtube.
Send a youtube url, select a format (mp3 or .mpeg) and you'll receive it.

# Getting Started

## Prerequisites
 - First of all you need a telegram bot api-key. (https://core.telegram.org/bots)
 - Docker installed If you want use the docker image
 - docker-compose installed if you want to use the compose file provided
 - ffmpeg installed if you want run the compiled Go binary

## Run Docker

If you want to generate the docker image I provide a Makefile (thanks to [rentziass](https://github.com/rentziass)).
You have to change the script and insert your api-key (maybe is better if is an env variable)

``` 
make
```

Now you can start the image with:

``` 
docker run yt-telegram-bot
```

Or with compose file:

``` 
docker-compose up .
```

If you want the container run in detached mode add **-d** flag.

## Run Binary
Build the binary on your own:
``` 
go build -o yt-telegram-bot
```
Then run it:
``` 
./yt-telegram-bot
```

# Libraries
Thanks to these libraries, they simplified my worklow:
 - [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
 - [kkdai/youtube](https://github.com/kkdai/youtube)
# License
distributed under MIT.


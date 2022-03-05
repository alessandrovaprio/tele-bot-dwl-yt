FROM alpine
RUN apk add  --no-cache ffmpeg
ADD yt-telegram-bot /
CMD ["/yt-telegram-bot"]
ARG TG_API_KEY
ENV TELEGRAM_API_KEY=$TG_API_KEY
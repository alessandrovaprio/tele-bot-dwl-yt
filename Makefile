test:
	go test ./...

build:
	CGO_ENABLED=0 go build -o yt-telegram-bot
	docker image build -t yt-telegram-bot --build-arg TG_API_KEY="{YOUR_API_KEY_HERE}" .
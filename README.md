# rochimfn/tika-bot

## Description
This bot will help you extract text from many kinds of formats (including Images).

## Quickstart
### Run
```bash
docker run --rm -d --name tika-bot -e TELEGRAM_TOKEN=xxx:xxx rochimfn/tika-bot:0.3
```

### Build and Run

```bash
git clone https://github.com/rochimfn/content-extract-bot.git
cd content-extract-bot
docker build -t rochimfn/tika-bot:0.3 . 
docker run --rm -d --name tika-bot -e TELEGRAM_TOKEN=xxx:xxx rochimfn/tika-bot:0.3
```
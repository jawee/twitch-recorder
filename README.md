# twitch-recorder
![example workflow](https://github.com/jawee/twitch-recorder/actions/workflows/build-and-test.yml/badge.svg)
# Work in progress. Barely configurable at all, and not very smart.

## Usage
Make a copy of config-example.json, named config.json and fill it with 
client-id, client-secret and streamers. 

If you want to get a discord notification on start recording, add webhook-id and webhook-token

```json
{
    "client-id": "asdfkölk93242340fdsf",
    "client-secret": "jkklajhdfhj88912313",
    "streamers": "streamer1, streamer2",
    "webhook-id": "asdfasdfaf",
    "webhook-token": "kjasdkjlfkasdjfkaljkf"
}
```


### Docker
Place config.json in your ./config directory.


#### CLI

```bash
$ docker create -v /path/to/config/directory:/config \ 
  -v /path/to/in-progress/directory:/inprogress \
  -v /path/to/finished/videos/directory:/videos \
  -v /path/to/logs:/logs \
  --restart unless-stopped \
  --name twitch-recorder ghcr.io/jawee/twitch-recorder:latest
```

#### Docker-compose

```yaml
---
version: "2.1"
services:
  twitch-recorder:
    image: ghcr.io/jawee/twitch-recorder:latest
    container_name: twitch-recorder
    volumes:
      - ./config:/config
      - ./inprogress:/inprogress
      - ./videos:/videos
      - ./logs:/logs
    restart: unless-stopped
```


### Development


#### Docker run

```bash
$ git clone git@github.com:jawee/twitch-recorder.git
$ cd twitch-recorder
$ docker build -t twitch-recorder .
$ docker run -it --rm -v /path/to/config/directory:/config \ 
-v /path/to/in-progress/directory:/inprogress \
-v /path/to/finished/videos/directory:/videos \
--name twitch-recorder twitch-recorder
```

#### Go run
```bash
go run cmd/main.go 
```

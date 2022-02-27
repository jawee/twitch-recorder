# twitch-recorder
![example workflow](https://github.com/jawee/twitch-recorder/actions/workflows/build-and-test.yml/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=jawee_twitch-recorder&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=jawee_twitch-recorder)
# Work in progress. Barely configurable at all, and not very smart.


## TODO
- [ ] Move finished files from inprogress to videos
- [ ] Convert videos to h265 with ffmpeg
- [ ] Notifications (probably discord, success and failure)

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

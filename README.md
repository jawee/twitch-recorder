## twitch-recorder

## Work in progress. Barely configurable at all, and not very smart.


#### TODO
- [ ] Move finished files from inprogress to videos
- [ ] Convert videos to h265 with ffmpeg

#### Usage
Make a copy of config-example.json, named config.json and fill it with 
client-id and client-secret from twitch. 


##### Docker

```bash
git clone git@github.com:jawee/twitch-recorder.git
cd twitch-recorder
docker build -t twitch-recorder .
docker run -it --rm -v /path/to/config/directory:/config \ 
-v /path/to/in-progress/directory:/inprogress \
-v /path/to/finished/videos/directory:/videos \
--name twitch-recorder twitch-recorder
```

##### Go run
```bash
go run cmd/main.go 
```

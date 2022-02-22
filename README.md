## twitch-recorder

#### Usage
Make a copy of config-example.json, named config.json and fill it with 
client-id and client-secret from twitch. 

##### Go run
```bash
go run cmd/main.go 
```

##### Docker

```bash
git clone git@github.com/jawee/twitch-recorder.git
cd twitch-recorder
docker build -t my-golang-app .
docker run -it --rm -v /path/to/in-progress/directory:/tempdir -v /path/to/finished/videos/directory:/videos --name my-running-app my-golang-app
```











##### Help

##### Streamlink usage
```bash
streamlink --twitch-disable-ads twitch.tv/username 1080p -o filetest.ts
```
```bash
streamlink --twitch-disable-ads twitch.tv/username best -o filetest.mp4/mkv/ts
```

How to detect which quality is best?

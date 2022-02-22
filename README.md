## twitch-recorder

## Work in progress. Barely configurable at all, and not very smart.

#### Usage
Make a copy of config-example.json, named config.json and fill it with 
client-id and client-secret from twitch. 

##### Go run
```bash
go run cmd/main.go 
```

##### Docker

```bash
git clone git@github.com:jawee/twitch-recorder.git
cd twitch-recorder
docker build -t my-golang-app .
docker run -it --rm -v /path/to/in-progress/directory:/tempdir -v /path/to/finished/videos/directory:/videos --name my-running-app my-golang-app
```

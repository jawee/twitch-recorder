FROM golang:1.17

WORKDIR /usr/src/app

RUN echo "deb http://deb.debian.org/debian buster-backports main" | tee "/etc/apt/sources.list.d/streamlink.list"

RUN apt update
RUN apt -t buster-backports install streamlink -y
# RUN apt install ffmpeg -y

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/...

RUN mkdir /inprogress
RUN mkdir /videos
RUN mkdir /config

CMD ["app"]

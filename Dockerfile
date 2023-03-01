FROM golang:1.20

WORKDIR /usr/src/app

RUN echo "deb http://deb.debian.org/debian bullseye-backports main" | tee "/etc/apt/sources.list.d/streamlink.list"

RUN apt update
RUN apt -t bullseye-backports install streamlink -y
# RUN apt install ffmpeg -y

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN echo "Files copied"
RUN go build -o /usr/local/bin/app -buildvcs=false ./cmd/... 

RUN mkdir /inprogress
RUN mkdir /videos
RUN mkdir /config
RUN mkdir /logs

CMD ["app"]

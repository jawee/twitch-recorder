FROM golang:1.24 AS build

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN echo "Files copied"
RUN go build -o /usr/local/bin/twitch-recorder -buildvcs=false ./cmd/... 


# Run the tests in the container
FROM build AS run-test-stage
RUN go test -v ./...


FROM debian:stable-slim AS release-stage

WORKDIR /

COPY --from=build /usr/local/bin/twitch-recorder /usr/local/bin/twitch-recorder

RUN echo "deb http://deb.debian.org/debian bookworm-backports main" | tee "/etc/apt/sources.list.d/streamlink.list"

RUN apt update
RUN apt -t bookworm-backports install streamlink -y

# RUN apt install ffmpeg -y

RUN mkdir /inprogress
RUN mkdir /videos
RUN mkdir /config
RUN mkdir /logs

CMD ["twitch-recorder"]

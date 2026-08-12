FROM golang:1.26 AS build

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN echo "Files copied"
RUN go build -o /usr/local/bin/twitch-recorder -buildvcs=false ./cmd/... 


# Run the tests in the container
FROM build AS run-test-stage
RUN go test -v ./...


FROM debian:trixie-slim AS release-stage

WORKDIR /

COPY --from=build /usr/local/bin/twitch-recorder /usr/local/bin/twitch-recorder
COPY docker/chromium-launcher /usr/local/bin/chromium-launcher

RUN echo "deb http://deb.debian.org/debian trixie-backports main" > /etc/apt/sources.list.d/streamlink.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        -t trixie-backports \
        streamlink \
    && apt-get install -y --no-install-recommends \
        chromium \
    && useradd --system --create-home --home-dir /var/lib/streamlink-browser streamlink-browser \
    && chmod 0755 /usr/local/bin/chromium-launcher \
    && rm -rf /var/lib/apt/lists/*

# RUN apt install ffmpeg -y

RUN mkdir /inprogress
RUN mkdir /videos
RUN mkdir /config
RUN mkdir /logs

CMD ["twitch-recorder"]

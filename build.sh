#! /usr/bin/env bash

if [ -e bin/playerMarkers* ]; then
    rm bin/playerMarkers*
fi

go build -o bin/playerMarkers-linux_amd64
env GOARCH=386 go build -o bin/playerMarkers-linux_i386

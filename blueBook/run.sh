#!/usr/bin/env bash
sudo docker -H localhost:2375 build --rm -t blue-book .
sudo docker -H localhost:2375 run -it --rm --name blue-book-1 blue-book
#!/bin/bash

mp3_dir="/var/www/html/helpersofyourjoy/media/mp3"
redis_prefix="HOYJ::MP3::MAP::DUMP"
for i in $( find "$mp3_dir" -iname "*.mp3" -print)
do
    http_path=$(echo "$i" | sed 's/\/var\/www\/html\/helpersofyourjoy/https:\/\/b.helpersofyourjoy.com/')
    echo "http_path: $http_path"
    basefile=$(basename "$i")
    echo "$basefile"
    redis-cli SADD "$redis_prefix" "$http_path"
done

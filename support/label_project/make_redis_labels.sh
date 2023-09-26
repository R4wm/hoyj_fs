#!/bin/bash

for f in $(find /var/www/html/helpersofyourjoy/storage/labels -iname "*json" -print)
do
    filename=$(basename $f)
    echo "filename: $filename"
    redis_label=$(cat $f | jq .md5sum | sed 's/\"//g')
    echo "redis_label: $redis_label"
    payload=$(cat $f | jq .)
    file_label="HOYJ::LABEL::$redis_label"
    echo "pushing redis: $file_label"
    redis-cli set  "$file_label" "$payload"
done


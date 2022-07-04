#!/bin/bash

TEMP="/tmp/list.md5.temp"
DEST="/tmp/list.md5"

echo "removing $TEMP"
rm -f "$TEMP"

for i in $(find /var/www/html/helpersofyourjoy/storage/media -type f -print )
do
    echo "working on $i"
    name=$(basename $i)
    md5result=$(md5sum "$i" | awk '{print $1}')
    echo "name: $name"
    echo "md5:  $md5result" 
    echo "$name $md5result" >> /tmp/md5
done

echo "Wrote to $DEST"

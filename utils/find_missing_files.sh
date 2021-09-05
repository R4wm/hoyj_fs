#!/bin/bash

for i in $(find . -type f -iname "*.mp3" -print)
do
    filename=$(basename $i)
    grep -q $filename /tmp/media.txt  || echo ">> not found in media.txt: $(realpath $i)"
done

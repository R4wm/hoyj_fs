#!/bin/bash
set -e
###############################
# Just downloading everything #
###############################
destination=~/media/hoyj/
for i in $(cat ./medialist.txt)
do
    subject=$(echo $i | cut -f1 -d'/')
    echo "subject: " $subject
    filename=$(echo $i | cut -f2 -d'/')
    echo "filename: " $filename
    if [ -f "$destination/$filename" ]
    then
        echo "skipping download, file exists: $destination/$filename"
    else
        linode-cli obj get "$subject" "$filename" "$destination/$filename"
    fi
done

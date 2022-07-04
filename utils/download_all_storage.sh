#!/bin/bash
set -e
###############################
# Just downloading everything #
###############################
destination=/var/www/html/helpersofyourjoy/storage/media
for i in $(cat ./medialist.txt)
do
    subject=$(echo $i | cut -f1 -d'/')
    echo "subject: " $subject
    if [ "$subject" = "collective" ]
    then
        echo "skipping collective"
        continue
    fi
    subDestination="$destination/$subject"
    mkdir -p "$subDestination"
    echo "Created $subDestination"
    
    filename=$(echo $i | cut -f2 -d'/')
    echo "filename: " $filename
    if [ -f "$subDestination/$filename" ]
    then
        echo "skipping download, file exists: $destination/$filename"
    else
        linode-cli obj get "$subject" "$filename" "$subDestination/$filename"
    fi
done

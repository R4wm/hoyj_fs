#!/bin/bash
# THis is a throw away script
# Get the directory names from Linode Object Storage and create those directories if not exist
# Use Linode directory as mapping to move files to appropriate directory locally


set -e
count=0

for i in $(cat ./medialist.txt)
do
    if [ $count -gt 10 ]
    then
        exit 0
    fi
    
    subject=$(echo $i | cut -f1 -d'/')
    echo "subject: " $subject
    if [ "$subject" = "collective" ]
    then
        echo "skipping collective"
        continue
    fi
    # make sure that directory exists
    destination=/var/www/html/helpersofyourjoy/storage/media/$subject
    mkdir -p $destination
    echo "directory: $destination"
    
    filename=$(echo $i | cut -f2 -d'/')
    echo "filename: " $filename
    if [ -f "$destination/$filename" ]
    then
        echo "skipping download, file exists: $destination/$filename"
    else
        echo "check me out: $destination/$filename"
        filepath=$(find /var/www/html/helpersofyourjoy/storage/media -name "$filename")
        mv "$filepath" "$destination"
        echo "moved $filepath -> $destination"
    fi
    echo "---------------------------------------------------------------"
    ((count=count+1))
done
11;rgb:0000/0000/0000

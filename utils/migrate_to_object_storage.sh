#!/bin/bash
set -e

lowercase_direcories(){
    for i in $(find . -type d -print)
    do
        bucket_name=$(echo $i | tr -d './' | tr '_' '-' | tr '[:upper:]' '[:lower:]')
        echo "this is bucket name: $bucket_name"
        linode-cli obj mb "$bucket_name"
    done
} 

# remove commas from file names
remove_commas(){
    echo removing commas
    for i in $(find . -iname "*,*" -print)
    do
        newname=$(echo $i | tr -d ',')
        echo "this is newname: $newname"
        mv -v $i $newname
    done

}

# remove apostrophes from filename
remove_apostrophe(){
    echo removing commas
    for i in $(find . -iname "*'*" -print)
    do
        echo "this is oldname: $i"
        newname=$(echo $i | tr -d "'")
        echo "this is newname: $newname"
        mv -v $i $newname
    done
    
}

# Migrate files over to Linode Object Storage
migrate(){
    echo "running migration"
    for i in $(find . \( -iname "*.mp3" -o -iname "*.mp4" \) -print)
    do
        _dirname=$(dirname "$i")
        _dirname=$(echo "$_dirname" | tr -d './' | tr '_' '-' | tr '[:upper:]' '[:lower:]')
        echo "original name: $i"
        echo "dirname: $_dirname"
        _filename=$(basename $i)
        echo "plain filename: $_filename"

        linode-cli obj put --acl-public "$i" "$_dirname"
    done
    
}


#remove_commas
#remove_apostrophe
migrate

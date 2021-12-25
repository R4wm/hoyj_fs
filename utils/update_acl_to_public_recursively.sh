#!/bin/bash
set -e
# Synopsis:
# Seems linode-cli doesnt offer a way to recursively set permissions from its help args
# Linode recommends using "s3" but aint nobody got time for that!
# ref: https://www.linode.com/community/questions/20759/how-can-i-change-all-existing-objects-in-my-bucket-to-public-or-private-using-s3
# script assumes "linode-cli" is properly configured

make_all_public_read(){
    echo "running make_all_public_read.. stop now if needed.."
    sleep 1
    for i in $(linode-cli obj la | awk '{print $4}')
    do
        _bucket=$(echo "$i" | cut -f1 -d/)
        _object=$(echo "$i" | cut -f2 -d/)
        echo "bucket: $_bucket"
        echo "object: $_object"
	linode-cli obj setacl --acl-public "$_bucket" "$_object"
    done
}

make_all_public_read


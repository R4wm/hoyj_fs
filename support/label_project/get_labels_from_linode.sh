#!/bin/bash
rm -f /tmp/media_list.txt
REDISKEY="HOYJ::LABEL::"

for i in $(linode-cli obj ls collective | awk '{print $4}')
do
    echo "$i" >> /tmp/media_list.txt
done



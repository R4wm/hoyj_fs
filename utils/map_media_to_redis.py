#!/usr/bin/env python3

# pull down media list from Linodes object storage
# create urls for each file
# push those urls to redis
# api knows how to parse search query and find url for related search query

import redis
from pathlib import Path

def populate_redis(a_media_file: str) -> None:
    base_url='us-southeast-1.linodeobjects.com'
    if not Path(a_media_file).is_file():
        print('{} does not exist!!'.format(a_media_file))
        return
    
    
    fileHandler = open(a_media_file, 'r')
    lines = fileHandler.readlines()
    redis_path = 'HOYJ::MEDIA'
    redis_handler = redis.Redis()
    
    for line in lines:
        complete_url = ''
        if line == '\n':
            continue
        else:
            line = line.split()[-1]
            message_title = line.split('/')[-1]
            url_prefix = line.split('/')[0]
            complete_url='https://{}.{}/{}'.format(url_prefix, base_url, message_title)
            #print('{}::{}'.format(redis_path, url_prefix.upper()))
        redis_handler.sadd('{}::{}'.format(redis_path, url_prefix.upper())  , complete_url)



populate_redis('/tmp/media.txt')

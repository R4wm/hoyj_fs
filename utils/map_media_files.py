#!/usr/bin/env python3

import redis
from pathlib import Path
def populate_redis(a_media_file: str) -> None:
    base_url='us-southeast-1.linodeobjects.com'
    if not Path(a_media_file).is_file():
        print('{} does not exist!!'.format(a_media_file))
        return
    
    
    fileHandler = open(a_media_file, 'r')
    lines = fileHandler.readlines()
    redis_path = 'HOYJ::MP3::MAP::DUMP'
    redis_inprogress = 'HOYJ::MP3::MAP::INPROGRESS'
    redis_handler = redis.Redis()
    
    count = 0
    for line in lines:
        complete_url = ''
        count += 1
        if line == '\n':
            print('>> line is empty')
            continue
        else:
            print("ok line: {}".format(line))
            line = line.split()[-1]
            message_title = line.split('/')[-1]
            print('this is message_title: {}'.format(message_title))
            print('line again: {}'.format(line))
            url_prefix = line.split('/')[0]
            print('url prefix: {}'.format(url_prefix))
            complete_url='https://{}.{}/{}'.format(url_prefix, base_url, message_title)
        print('this is complete url: {}'.format(complete_url))
        redis_handler.sadd(redis_inprogress , complete_url)

    redis_handler.rename(redis_inprogress, redis_path)

populate_redis('/tmp/media.txt')

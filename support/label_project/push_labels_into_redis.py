#!/usr/bin/env python3
import redis
import requests
import json
import walk

r = redis.Redis(host='localhost', port=6379)
REDISPREFIX = "HOYJ::LABEL::"
PREFIX='https://collective.us-southeast-1.linodeobjects.com/'
LABELPATH='/var/www/html/helpersofyourjoy/storage/labels'
JSONFILES = []

# ls the label directory, get all json
# read each json and dump contents into redis

for (dirpath, dirnames, filenames) in walk(mypath):
    f.extend(filenames)
    break

# with open('/tmp/media_list.txt', 'r') as a:
#     for line in a:
#         # check if the key in redis already exists, skip it true
#         redis_key = REDISPREFIX + line.strip()
#         check_redis = r.exists(redis_key)
#         if check_redis > 0:
#             print("skipping ", redis_key)
#             continue
#         # import pdb;pdb.set_trace()
#         url = PREFIX + line.strip()
#         print(url)
#         result = requests.get(url)
#         if result.status_code != 200:
#             FAILED.append(url)
#             continue
#         print(result.text)
#         jsonified = json.loads(result.text)
#         print("redis key: ", redis_key)
#         r.set(redis_key, result.text)

        
        
# print("FAILED: ")
# print(FAILED)

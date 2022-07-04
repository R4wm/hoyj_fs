#!/usr/bin/env python3

import os
import redis
import json


redisKey = 'LABELS'
r = redis.Redis()
PATH='/var/www/html/helpersofyourjoy/storage/labels/'
ls = os.listdir(PATH)

for i in ls:
    fullPath=PATH+i
    with open(fullPath) as e:
        payload = json.loads(e.read())
        print(payload)
        r.sadd(redisKey, json.dumps(payload))

#!/bin/bash

linode-cli obj la > /tmp/media.txt && ./map_media_to_redis.py && ./map_media_files.py

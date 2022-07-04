#!/usr/bin/env python3

import os
import json
import sqlite3

con = sqlite3.connect('/tmp/labels.db')
cur = con.cursor()

cur.execute('''
CREATE TABLE IF NOT EXISTS labels (
originalref text,
filename text,
speaker text,
year int,
month int,
day int,
topic text,
md5sum text,
book text,
chapter int,
part int,
verseStart int,
verseEnd int,
outlinefile text,
URL text
)
''')

con.commit()


PATH='/var/www/html/helpersofyourjoy/storage/labels/'
ls = os.listdir(PATH)

for i in ls:
    fullPath=PATH+i
    with open(fullPath) as e:
        payload = json.loads(e.read())
        print(payload)
        cur.execute("INSERT INTO labels(originalref, filename, speaker, year, month, day, topic, md5sum, book, chapter, part, verseStart, verseEnd, outlinefile, URL) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",(
                    payload["originalref"],
                    payload["filename"],
                    payload["speaker"],
                    payload["year"],
                    payload["month"],
                    payload["day"],
                    payload["topic"],
                    payload["md5sum"],
                    payload["book"],
                    payload["chapter"],
                    payload["part"],
                    payload["verseStart"],
                    payload["verseEnd"],
                    payload["outlinefile"],
                    payload["URL"],)
        )
        con.commit()


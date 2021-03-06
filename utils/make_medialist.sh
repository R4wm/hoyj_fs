#!/bin/bash
TEMP="/tmp/media.pending.txt"
DEST="/tmp/media.txt"
echo "rm $TEMP"

rm -f "$TEMP"

for i in $(linode-cli obj ls | awk '{print $3}')
do
    echo "working on $i"
    for j in $(linode-cli obj ls $i | awk '{print $4}')
    do
        echo "$i/$j" >> "$TEMP"
    done
done

echo "mv $TEMP -> $DEST"
rm -f "$DEST"
mv "$TEMP" "$DEST" 
echo "wrote media list to $DEST"


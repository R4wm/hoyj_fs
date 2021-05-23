#!/bin/bash
COMEBACK=$(pwd)

###############################
# REMOVE ALL THE WHITE SPACES #
###############################
# find . -type d | while read dir
# do
#     cd "$dir"
#     echo "you are here: $(pwd)"
#     for f in *
#     do
# 	new="${f// /_}"
# 	if [ "$new" != "$f" ]
# 	then
# 	    if [ -e "$new" ]
# 	    then
# 		echo not renaming \""$f"\" because \""$new"\" already exists
# 	    else
# 		echo moving "$f" to "$new"
# 		mv "$f" "$new"
# 	    fi
# 	fi
#     done
#     cd $COMEBACK
# done


#####################################
# remove all files with "*.mp3.mp3" #
#####################################
# for i in $(find . -iname "*.mp3.mp3" -print)
# do
#     echo "old: $i"
#     newfilename="$(echo $i | sed 's/.mp3.mp3/.mp3/')"
#     echo "new: $newfilename"
#     mv "$i" "$newfilename"
# done



##############################
# remove double periods (..) #
##############################
# echo "removing double periods"
# for i in $(find . -iname "*..*" -print) 
# do
#     echo "old: $i"
#     newfilename="$(echo $i | sed 's/\.\./\./')"
#     echo "new: $newfilename"
#     mv "$i" "$newfilename"
# done

# echo "finished removing double period"

###################################
# Uniform date format to YYYYMMDD #
###################################
# for i in $(find . -iname "20[0-9][0-9]-*.mp3" -print)
# do
#     echo ""
#     echo "original: $i"
#     basefile=$(basename "$i")
#     echo "old: $basefile"
#     year=$(echo $basefile | awk '{print substr($0,1,4)}')
#     echo "year: $year"

#     month=$(echo $basefile | awk '{print substr($0,6,2)}')
#     echo "month: $month"

#     day=$(echo $basefile | awk '{print substr($0,9,2)}')
#     echo "day: $day"
    
#     remainder=$(echo $basefile | awk '{print substr($0,11)}')
#     echo "remainder: $remainder"
    
    
#     postname=$(echo $basefile | cut -f3 -d'-')
#     echo "prename: $postname"

#     parentdir=$(dirname "$i")
#     newfilename="$parentdir/$year$month$day$remainder"
#     echo "newfilename: $newfilename"

#     mv "$i" "$newfilename"

# done


################
# TAG Speakers #
################
for i in $(find . -iname "*richard*jordan*" -print)
do
    echo "$i"
    kid3-cli -c "set artist richard_jordan"
    echo "finished tagging speaker $i"
done


for i in $(find . -iname "*nate*cody*" -print)
do
    echo "$i"
    kid3-cli -c "set artist nate_cody"
    echo "finished tagging speaker $i"
done

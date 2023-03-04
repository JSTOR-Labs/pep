#!/bin/sh
BASEDIR=$(dirname "$0")
osascript -e "tell application \"Terminal\"" -e "do script \"$BASEDIR/JSTOR/elasticsearch/bin/elasticsearch\"" -e "end tell"
sleep 60
$BASEDIR/JSTOR/api-mac serve
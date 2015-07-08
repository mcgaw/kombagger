#!/bin/bash

# Watch all *.go files in the specified directory
# Call the restart function when they are saved
function monitor() {
  inotifywait -q -m -r @/home/mcgaw/kombagger/src/github.com/mcgaw/kombagger/ember-kombagger -e close_write -e moved_to --exclude '[^g][^o]$' $1 |
  while read line; do
    restart
  done
}

# Terminate and rerun the main Go program
function restart {
  if [ "$(pidof $PROCESS_NAME)" ]; then
    killall -q -w -9 $PROCESS_NAME
  fi
  echo ">> Reloading..."
  go run $FILE_PATH $ARGS &
}

# Make sure all background processes get terminated
function close {
  killall -q -w -9 inotifywait
  exit 0
}

trap close INT
echo "== Go-reload"
echo ">> Watching directories, CTRL+C to stop"

FILE_PATH=$1
FILE_NAME=$(basename $FILE_PATH)
PROCESS_NAME=${FILE_NAME%%.*}

shift
ARGS=$@

# Start the main Go program
go run $FILE_PATH $ARGS &

# Monitor all /src directories on the GOPATH
OIFS="$IFS"
IFS=':'
for path in $GOPATH
do
  monitor $path/src &
done
IFS="$OIFS"

# If the current working directory isn't on the GOPATH, monitor it too
if [[ $PWD != "$GOPATH/"* ]]
then
  monitor $PWD
fi

wait

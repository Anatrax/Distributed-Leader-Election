#!/bin/bash
rm output.txt
pkill -f better_grpc_base

go build

x=$1
if [ $# -eq 0 ]; then
    x=10
fi

for ((i=0;i<$x;i++))
do
    echo "Starting node ($i/$x)..."
	./better_grpc_base $x $i >> output.txt &
    # sleep 1
done

tail -f output.txt &

read quit
./stop.sh
